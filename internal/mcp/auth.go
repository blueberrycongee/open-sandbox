package mcp

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	Enabled      bool
	JWTSecret    string
	JWTPublicKey string
	Audience     string
	Issuer       string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		Enabled:      readBoolEnv("MCP_AUTH_ENABLED"),
		JWTSecret:    strings.TrimSpace(os.Getenv("MCP_AUTH_JWT_SECRET")),
		JWTPublicKey: strings.TrimSpace(os.Getenv("MCP_AUTH_JWT_PUBLIC_KEY")),
		Audience:     strings.TrimSpace(os.Getenv("MCP_AUTH_AUDIENCE")),
		Issuer:       strings.TrimSpace(os.Getenv("MCP_AUTH_ISSUER")),
	}
}

func (config AuthConfig) Validate() error {
	if !config.Enabled {
		return nil
	}
	if config.JWTSecret == "" && config.JWTPublicKey == "" {
		return errors.New("jwt secret or public key required")
	}
	if config.JWTSecret != "" && config.JWTPublicKey != "" {
		return errors.New("only one jwt key source is allowed")
	}
	return nil
}

type Authenticator struct {
	config  AuthConfig
	key     any
	methods []string
	parser  *jwt.Parser
}

func NewAuthenticator(config AuthConfig) (*Authenticator, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	auth := &Authenticator{
		config: config,
		parser: jwt.NewParser(),
	}
	if !config.Enabled {
		return auth, nil
	}

	if config.JWTSecret != "" {
		auth.key = []byte(config.JWTSecret)
		auth.methods = []string{"HS256", "HS384", "HS512"}
		return auth, nil
	}

	key, methods, err := parsePublicKey(config.JWTPublicKey)
	if err != nil {
		return nil, err
	}
	auth.key = key
	auth.methods = methods
	return auth, nil
}

func (auth *Authenticator) Enabled() bool {
	if auth == nil {
		return false
	}
	return auth.config.Enabled
}

func (auth *Authenticator) ValidateRequest(r *http.Request) *ErrorDetail {
	if auth == nil || !auth.config.Enabled {
		return nil
	}
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return unauthorizedError("missing bearer token")
	}
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return unauthorizedError("invalid authorization header")
	}
	if parts[1] == "" {
		return unauthorizedError("missing bearer token")
	}

	claims := jwt.MapClaims{}
	_, err := auth.parser.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
		if len(auth.methods) > 0 {
			for _, method := range auth.methods {
				if token.Method.Alg() == method {
					return auth.key, nil
				}
			}
			return nil, errors.New("unexpected signing method")
		}
		return auth.key, nil
	})
	if err != nil {
		return unauthorizedError("invalid token")
	}
	if auth.config.Audience != "" && !audienceMatches(claims["aud"], auth.config.Audience) {
		return unauthorizedError("invalid audience")
	}
	if auth.config.Issuer != "" {
		issuer, _ := claims["iss"].(string)
		if issuer != auth.config.Issuer {
			return unauthorizedError("invalid issuer")
		}
	}
	return nil
}

func unauthorizedError(message string) *ErrorDetail {
	detail := NewErrorDetail("unauthorized", message, KindUnauthorized)
	return &detail
}

func parsePublicKey(raw string) (any, []string, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(raw))
	if err == nil {
		return key, []string{"RS256", "RS384", "RS512"}, nil
	}
	if ecdsaKey, ecdsaErr := jwt.ParseECPublicKeyFromPEM([]byte(raw)); ecdsaErr == nil {
		return ecdsaKey, []string{"ES256", "ES384", "ES512"}, nil
	}
	if edKey, edErr := jwt.ParseEdPublicKeyFromPEM([]byte(raw)); edErr == nil {
		return edKey, []string{"EdDSA"}, nil
	}
	return nil, nil, errors.New("unsupported public key format")
}

func audienceMatches(raw any, audience string) bool {
	switch value := raw.(type) {
	case string:
		return value == audience
	case []string:
		for _, entry := range value {
			if entry == audience {
				return true
			}
		}
	case []any:
		for _, entry := range value {
			if str, ok := entry.(string); ok && str == audience {
				return true
			}
		}
	}
	return false
}

func readBoolEnv(key string) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return false
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

var _ = rsa.PublicKey{}
var _ = ecdsa.PublicKey{}
var _ = ed25519.PublicKey{}
