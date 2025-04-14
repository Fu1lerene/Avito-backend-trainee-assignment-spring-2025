package errors

import (
	core_errors "avito/internal/models/error"
	"encoding/json"
	"errors"
	"net/http"
)

func Write(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	var appErr *core_errors.Error
	if errors.As(err, &appErr) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(appErr)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(core_errors.Error{
		Message: "unexpected error",
	})
}

var (
	ErrInvalidRole        = core_errors.NewError("invalid role")
	ErrInvalidToken       = core_errors.NewError("invalid token")
	ErrInvalidRequest     = core_errors.NewError("invalid request body")
	ErrUserNotFound       = core_errors.NewError("user not found")
	ErrInvalidCredentials = core_errors.NewError("invalid credentials")
	ErrAccessDenied       = core_errors.NewError("access denied")
	ErrInvalidProductType = core_errors.NewError("invalid product type")
	ErrNoActiveReceptions = core_errors.NewError("no active receptions")
	ErrNoProductsToDelete = core_errors.NewError("no products to delete")
	ErrDeleteReception    = core_errors.NewError("reception already closed or does not exist")
	ErrInvalidCity        = core_errors.NewError("city not allowed")
	ErrReceptionNotClosed = core_errors.NewError("reception not closed")
)
