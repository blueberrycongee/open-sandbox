package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Server struct {
	registry *Registry
	auth     *Authenticator
	authErr  error
}

func NewServer(registry *Registry, auth *Authenticator, authErr error) *Server {
	return &Server{
		registry: registry,
		auth:     auth,
		authErr:  authErr,
	}
}

func (server *Server) HandleRequest(ctx context.Context, req Request) Response {
	if req.JSONRPC != JSONRPCVersion {
		detail := NewErrorDetail(KindInvalidRequest, "invalid jsonrpc version", KindInvalidRequest)
		return NewErrorResponse(req.ID, ErrInvalidRequest, "invalid request", detail)
	}
	if err := ValidateProtocolVersion(req.Params); err != nil {
		detail := NewErrorDetail(KindInvalidParams, err.Error(), KindInvalidParams)
		return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
	}
	if req.Method == "mcp.capabilities" {
		return NewSuccessResponse(req.ID, BuildCapabilities(server.registry))
	}
	tool, ok := server.registry.Get(req.Method)
	if !ok || tool.Handler == nil {
		detail := NewErrorDetail(KindMethodNotFound, "unknown method", KindMethodNotFound)
		return NewErrorResponse(req.ID, ErrMethodNotFound, "method not found", detail)
	}
	result, toolErr := tool.Handler(ctx, req.Params)
	if toolErr != nil {
		if toolErr.TraceID == "" {
			toolErr.TraceID = NewErrorDetail(toolErr.Code, toolErr.Message, toolErr.Kind).TraceID
		}
		switch toolErr.Kind {
		case KindInvalidParams:
			return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", *toolErr)
		case KindUnauthorized:
			return NewErrorResponse(req.ID, ErrUnauthorized, "unauthorized", *toolErr)
		case KindForbidden:
			return NewErrorResponse(req.ID, ErrForbidden, "forbidden", *toolErr)
		default:
			return NewErrorResponse(req.ID, ErrToolExecution, "tool error", *toolErr)
		}
	}
	return NewSuccessResponse(req.ID, result)
}

func (server *Server) ServeStdio(r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(bufio.NewReader(r))
	encoder := json.NewEncoder(w)

	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		var req Request
		if err := json.Unmarshal(raw, &req); err != nil {
			detail := NewErrorDetail(KindInvalidRequest, "invalid request", KindInvalidRequest)
			resp := NewErrorResponse(nil, ErrInvalidRequest, "invalid request", detail)
			if err := encoder.Encode(resp); err != nil {
				return err
			}
			continue
		}

		parsed, err := ParseRequest(raw)
		if err != nil {
			detail := NewErrorDetail(KindInvalidRequest, err.Error(), KindInvalidRequest)
			resp := NewErrorResponse(req.ID, ErrInvalidRequest, "invalid request", detail)
			if err := encoder.Encode(resp); err != nil {
				return err
			}
			continue
		}

		resp := server.HandleRequest(context.Background(), parsed)
		if isNotification(parsed.ID) {
			continue
		}
		if err := encoder.Encode(resp); err != nil {
			return err
		}
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, notify := server.handleHTTPRequest(r)
	if notify {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (server *Server) ServeSSE(w http.ResponseWriter, r *http.Request) {
	resp, notify := server.handleSSERequest(r)
	if notify {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	payload, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("event: message\n"))
	_, _ = w.Write([]byte("data: "))
	_, _ = w.Write(payload)
	_, _ = w.Write([]byte("\n\n"))
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (server *Server) handleHTTPRequest(r *http.Request) (Response, bool) {
	if server.authErr != nil {
		detail := NewErrorDetail(KindInternal, server.authErr.Error(), KindInternal)
		return NewErrorResponse(nil, ErrInternal, "internal error", detail), false
	}
	if server.auth != nil {
		if authErr := server.auth.ValidateRequest(r); authErr != nil {
			return NewErrorResponse(nil, ErrUnauthorized, "unauthorized", *authErr), false
		}
	}
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		detail := NewErrorDetail(KindInvalidRequest, "unable to read request", KindInvalidRequest)
		return NewErrorResponse(nil, ErrInvalidRequest, "invalid request", detail), false
	}
	req, err := ParseRequest(payload)
	if err != nil {
		detail := NewErrorDetail(KindInvalidRequest, err.Error(), KindInvalidRequest)
		return NewErrorResponse(nil, ErrInvalidRequest, "invalid request", detail), false
	}
	resp := server.HandleRequest(r.Context(), req)
	return resp, isNotification(req.ID)
}

func (server *Server) handleSSERequest(r *http.Request) (Response, bool) {
	if server.authErr != nil {
		detail := NewErrorDetail(KindInternal, server.authErr.Error(), KindInternal)
		return NewErrorResponse(nil, ErrInternal, "internal error", detail), false
	}
	if server.auth != nil {
		if authErr := server.auth.ValidateRequest(r); authErr != nil {
			return NewErrorResponse(nil, ErrUnauthorized, "unauthorized", *authErr), false
		}
	}
	payload := r.URL.Query().Get("request")
	if payload == "" {
		detail := NewErrorDetail(KindInvalidRequest, "missing request", KindInvalidRequest)
		return NewErrorResponse(nil, ErrInvalidRequest, "invalid request", detail), false
	}
	req, err := ParseRequest([]byte(payload))
	if err != nil {
		detail := NewErrorDetail(KindInvalidRequest, err.Error(), KindInvalidRequest)
		return NewErrorResponse(nil, ErrInvalidRequest, "invalid request", detail), false
	}
	resp := server.HandleRequest(r.Context(), req)
	return resp, isNotification(req.ID)
}

func isNotification(id json.RawMessage) bool {
	if len(id) == 0 {
		return true
	}
	return bytes.Equal(bytes.TrimSpace(id), []byte("null"))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
