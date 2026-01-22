package types

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

type Response struct {
	Status string       `json:"status"`
	Data   any          `json:"data,omitempty"`
	Error  *ErrorDetail `json:"error,omitempty"`
}

func Ok(data any) Response {
	return Response{Status: StatusOK, Data: data}
}

func Fail(err *ErrorDetail) Response {
	return Response{Status: StatusError, Error: err}
}
