package mcp

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

const (
	JSONRPCVersion           = "2.0"
	SupportedProtocolVersion = "1.0"
	MethodCapabilities       = "mcp.capabilities"
	MethodInitialize         = "initialize"
	MethodToolsList          = "tools/list"
	MethodToolsCall          = "tools/call"
	ErrInvalidRequest        = -32600
	ErrMethodNotFound        = -32601
	ErrInvalidParams         = -32602
	ErrInternal              = -32603
	ErrUnauthorized          = -32001
	ErrForbidden             = -32003
	ErrToolExecution         = -32010
	KindInvalidRequest       = "invalid_request"
	KindMethodNotFound       = "method_not_found"
	KindInvalidParams        = "invalid_params"
	KindInternal             = "internal"
	KindUnauthorized         = "unauthorized"
	KindForbidden            = "forbidden"
	KindToolError            = "tool_error"
)

type JSONSchema map[string]any

type ToolSchema struct {
	Input  JSONSchema `json:"input,omitempty"`
	Output JSONSchema `json:"output,omitempty"`
}

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *ErrorDetail `json:"data,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
	Kind    string `json:"kind,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string `json:"protocol_version,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string `json:"protocol_version"`
}

type ToolsListResult struct {
	Tools []ToolInfo `json:"tools"`
}

type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

type ToolCallResult struct {
	Result any `json:"result"`
}

func ParseRequest(payload []byte) (Request, error) {
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return Request{}, err
	}
	if req.JSONRPC != JSONRPCVersion {
		return Request{}, errors.New("invalid jsonrpc version")
	}
	if req.Method == "" {
		return Request{}, errors.New("method is required")
	}
	return req, nil
}

func ValidateProtocolVersion(params json.RawMessage) error {
	if len(params) == 0 {
		return nil
	}
	var payload map[string]any
	if err := json.Unmarshal(params, &payload); err != nil {
		return nil
	}
	raw, ok := payload["protocol_version"]
	if !ok {
		return nil
	}
	version, ok := raw.(string)
	if !ok || strings.TrimSpace(version) == "" {
		return errors.New("invalid protocol version")
	}
	if version != SupportedProtocolVersion {
		return errors.New("unsupported protocol version")
	}
	return nil
}

func NewSuccessResponse(id json.RawMessage, result any) Response {
	return Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

func NewErrorDetail(code, message, kind string) ErrorDetail {
	return ErrorDetail{
		Code:    code,
		Message: message,
		TraceID: newTraceID(),
		Kind:    kind,
	}
}

func NewInvalidRequestDetail(message string) ErrorDetail {
	return NewErrorDetail(KindInvalidRequest, message, KindInvalidRequest)
}

func NewInvalidParamsDetail(message string) ErrorDetail {
	return NewErrorDetail(KindInvalidParams, message, KindInvalidParams)
}

func NewMethodNotFoundDetail(message string) ErrorDetail {
	return NewErrorDetail(KindMethodNotFound, message, KindMethodNotFound)
}

func NewInternalDetail(message string) ErrorDetail {
	return NewErrorDetail(KindInternal, message, KindInternal)
}

func NewErrorResponse(id json.RawMessage, code int, message string, detail ErrorDetail) Response {
	return Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
			Data:    &detail,
		},
	}
}

func newTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
