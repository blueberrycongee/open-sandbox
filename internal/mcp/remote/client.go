package remote

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"open-sandbox/internal/mcp"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (client *Client) ToolsList(ctx context.Context, cfg ServerConfig) (mcp.ToolsListResult, error) {
	req := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      []byte("1"),
		Method:  mcp.MethodToolsList,
		Params:  nil,
	}
	resp, err := client.doRequest(ctx, cfg, req)
	if err != nil {
		return mcp.ToolsListResult{}, err
	}
	if resp.Error != nil {
		return mcp.ToolsListResult{}, errors.New(resp.Error.Message)
	}
	raw, err := json.Marshal(resp.Result)
	if err != nil {
		return mcp.ToolsListResult{}, err
	}
	var result mcp.ToolsListResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return mcp.ToolsListResult{}, err
	}
	return result, nil
}

func (client *Client) ToolsCall(ctx context.Context, cfg ServerConfig, name string, args json.RawMessage) (mcp.ToolCallResult, error) {
	params := mcp.ToolCallParams{
		Name:      name,
		Arguments: args,
	}
	rawParams, err := json.Marshal(params)
	if err != nil {
		return mcp.ToolCallResult{}, err
	}
	req := mcp.Request{
		JSONRPC: mcp.JSONRPCVersion,
		ID:      []byte("1"),
		Method:  mcp.MethodToolsCall,
		Params:  rawParams,
	}
	resp, err := client.doRequest(ctx, cfg, req)
	if err != nil {
		return mcp.ToolCallResult{}, err
	}
	if resp.Error != nil {
		return mcp.ToolCallResult{}, errors.New(resp.Error.Message)
	}
	raw, err := json.Marshal(resp.Result)
	if err != nil {
		return mcp.ToolCallResult{}, err
	}
	var result mcp.ToolCallResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return mcp.ToolCallResult{}, err
	}
	return result, nil
}

func (client *Client) doRequest(ctx context.Context, cfg ServerConfig, req mcp.Request) (mcp.Response, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return mcp.Response{}, err
	}
	transport := strings.ToLower(strings.TrimSpace(cfg.Transport))
	if transport == "" || transport == "http" {
		return client.doHTTP(ctx, cfg, payload)
	}
	if transport == "sse" {
		return client.doSSE(ctx, cfg, payload)
	}
	return mcp.Response{}, fmt.Errorf("unsupported transport: %s", cfg.Transport)
}

func (client *Client) doHTTP(ctx context.Context, cfg ServerConfig, payload []byte) (mcp.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return mcp.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	applyHeaders(req, cfg)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return mcp.Response{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.Response{}, err
	}
	var parsed mcp.Response
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return mcp.Response{}, err
	}
	return parsed, nil
}

func (client *Client) doSSE(ctx context.Context, cfg ServerConfig, payload []byte) (mcp.Response, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return mcp.Response{}, err
	}
	query := u.Query()
	query.Set("request", string(payload))
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return mcp.Response{}, err
	}
	req.Header.Set("Accept", "text/event-stream")
	applyHeaders(req, cfg)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return mcp.Response{}, err
	}
	defer resp.Body.Close()
	message, err := readSSEMessage(resp.Body)
	if err != nil {
		return mcp.Response{}, err
	}
	var parsed mcp.Response
	if err := json.Unmarshal(message, &parsed); err != nil {
		return mcp.Response{}, err
	}
	return parsed, nil
}

func applyHeaders(req *http.Request, cfg ServerConfig) {
	if cfg.AuthorizationToken != "" && !hasAuthorizationHeader(cfg.Headers) {
		req.Header.Set("Authorization", "Bearer "+cfg.AuthorizationToken)
	}
	for key, value := range cfg.Headers {
		if strings.TrimSpace(key) == "" {
			continue
		}
		req.Header.Set(key, value)
	}
}

func hasAuthorizationHeader(headers map[string]string) bool {
	for key := range headers {
		if strings.EqualFold(key, "Authorization") {
			return true
		}
	}
	return false
}

func readSSEMessage(reader io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(reader)
	var data bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			data.WriteString(strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if data.Len() == 0 {
		return nil, errors.New("empty sse response")
	}
	return data.Bytes(), nil
}
