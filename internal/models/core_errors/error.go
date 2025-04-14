package core_errors

import "errors"

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrUserNotFound       = errors.New("user not found")
	ErrAccessDenied       = errors.New("access denied")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidProductType = errors.New("invalid product type")
	ErrNoActiveReceptions = errors.New("no active receptions")
	ErrNoProductsToDelete = errors.New("no products to delete")
	ErrDeleteReception    = errors.New("reception already closed or does not exist")
	ErrInvalidCity        = errors.New("city not allowed")
	ErrReceptionNotClosed = errors.New("reception not closed")
)
