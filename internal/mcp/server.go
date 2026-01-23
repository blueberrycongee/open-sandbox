package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
	switch req.Method {
	case MethodInitialize:
		return server.handleInitialize(req)
	case MethodToolsList:
		return server.handleToolsList(req)
	case MethodToolsCall:
		return server.handleToolsCall(ctx, req)
	case MethodCapabilities:
		return NewSuccessResponse(req.ID, BuildCapabilities(server.registry))
	}
	tool, ok := server.registry.Get(req.Method)
	if !ok || tool.Handler == nil {
		detail := NewErrorDetail(KindMethodNotFound, "unknown method", KindMethodNotFound)
		return NewErrorResponse(req.ID, ErrMethodNotFound, "method not found", detail)
	}
	return server.handleToolInvocation(ctx, req.ID, tool, req.Params, false)
}

func (server *Server) handleInitialize(req Request) Response {
	version := SupportedProtocolVersion
	if parsed, ok, err := extractProtocolVersion(req.Params); err != nil {
		detail := NewInvalidParamsDetail("invalid params")
		return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
	} else if ok {
		version = parsed
	}
	return NewSuccessResponse(req.ID, InitializeResult{
		ProtocolVersion: version,
		Capabilities: InitializeCapabilities{
			Tools: &InitializeToolsCapabilities{
				ListChanged: false,
			},
		},
		ServerInfo: InitializeServerInfo{
			Name:    ServerName,
			Version: ServerVersion,
		},
	})
}

func (server *Server) handleToolsList(req Request) Response {
	if len(req.Params) > 0 {
		var payload map[string]any
		if err := json.Unmarshal(req.Params, &payload); err != nil {
			detail := NewInvalidParamsDetail("invalid params")
			return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
		}
	}
	return NewSuccessResponse(req.ID, ToolsListResult{Tools: server.registry.List()})
}

func (server *Server) handleToolsCall(ctx context.Context, req Request) Response {
	if len(req.Params) == 0 {
		detail := NewInvalidParamsDetail("invalid params")
		return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
	}
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		detail := NewInvalidParamsDetail("invalid params")
		return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
	}
	if params.Name == "" {
		detail := NewInvalidParamsDetail("name is required")
		return NewErrorResponse(req.ID, ErrInvalidParams, "invalid params", detail)
	}
	tool, ok := server.registry.Get(params.Name)
	if !ok || tool.Handler == nil {
		detail := NewMethodNotFoundDetail("unknown tool")
		return NewErrorResponse(req.ID, ErrMethodNotFound, "method not found", detail)
	}
	return server.handleToolInvocation(ctx, req.ID, tool, params.Arguments, true)
}

func (server *Server) handleToolInvocation(ctx context.Context, id json.RawMessage, tool Tool, params json.RawMessage, wrapResult bool) Response {
	result, toolErr := tool.Handler(ctx, params)
	if toolErr != nil {
		return toolErrorResponse(id, toolErr)
	}
	if wrapResult {
		return NewSuccessResponse(id, newToolCallResult(result))
	}
	return NewSuccessResponse(id, result)
}

func toolErrorResponse(id json.RawMessage, toolErr *ErrorDetail) Response {
	if toolErr.TraceID == "" {
		toolErr.TraceID = NewErrorDetail(toolErr.Code, toolErr.Message, toolErr.Kind).TraceID
	}
	switch toolErr.Kind {
	case KindInvalidParams:
		return NewErrorResponse(id, ErrInvalidParams, "invalid params", *toolErr)
	case KindUnauthorized:
		return NewErrorResponse(id, ErrUnauthorized, "unauthorized", *toolErr)
	case KindForbidden:
		return NewErrorResponse(id, ErrForbidden, "forbidden", *toolErr)
	default:
		return NewErrorResponse(id, ErrToolExecution, "tool error", *toolErr)
	}
}

func newToolCallResult(result any) ToolCallResult {
	content := []ContentBlock{}
	if result != nil {
		if text, ok := result.(string); ok {
			content = append(content, ContentBlock{Type: "text", Text: text})
		} else if payload, err := json.Marshal(result); err == nil {
			content = append(content, ContentBlock{Type: "text", Text: string(payload)})
		}
	}
	return ToolCallResult{
		Content:           content,
		StructuredContent: result,
		Result:            result,
	}
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
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if wantsEventStream(r) {
		writeSSE(w, resp)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (server *Server) ServeSSE(w http.ResponseWriter, r *http.Request) {
	resp, notify := server.handleSSERequest(r)
	if notify {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	writeSSE(w, resp)
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

func writeSSE(w http.ResponseWriter, resp Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	writeSSEMessage(w, resp)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func writeSSEMessage(w io.Writer, resp Response) {
	payload, _ := json.Marshal(resp)
	_, _ = w.Write([]byte("event: message\n"))
	_, _ = w.Write([]byte("data: "))
	_, _ = w.Write(payload)
	_, _ = w.Write([]byte("\n\n"))
}

func wantsEventStream(r *http.Request) bool {
	accept := strings.ToLower(r.Header.Get("Accept"))
	return strings.Contains(accept, "text/event-stream")
}
