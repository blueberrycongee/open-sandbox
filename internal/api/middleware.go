package api

import (
	"encoding/json"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) *AppError

func WrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				WriteAppError(w, NewAppError(CodeInternalError, "internal error", http.StatusInternalServerError))
			}
		}()

		if err := handler(w, r); err != nil {
			WriteAppError(w, err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteAppError(w http.ResponseWriter, appErr *AppError) {
	if appErr == nil {
		appErr = NewAppError(CodeInternalError, "internal error", http.StatusInternalServerError)
	}
	if appErr.TraceID == "" {
		appErr.TraceID = NewTraceID()
	}
	_ = WriteJSON(w, appErr.HTTPStatus, ErrorResponse(appErr))
}
