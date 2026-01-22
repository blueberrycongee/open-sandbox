package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"open-sandbox/pkg/types"
)

const (
	CodeNotFound         = "not_found"
	CodeMethodNotAllowed = "method_not_allowed"
	CodeInternalError    = "internal_error"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	TraceID    string
}

func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func NewTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}

func WithTraceID(err *AppError, traceID string) *AppError {
	if err == nil {
		return &AppError{
			Code:       CodeInternalError,
			Message:    "internal error",
			HTTPStatus: http.StatusInternalServerError,
			TraceID:    traceID,
		}
	}
	err.TraceID = traceID
	return err
}

func ErrorResponse(err *AppError) types.Response {
	if err == nil {
		err = &AppError{
			Code:       CodeInternalError,
			Message:    "internal error",
			HTTPStatus: http.StatusInternalServerError,
		}
	}
	return types.Fail(&types.ErrorDetail{
		Code:    err.Code,
		Message: err.Message,
		TraceID: err.TraceID,
	})
}
