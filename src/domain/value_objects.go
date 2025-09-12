package domain

import (
	"errors"
)

var (
	ErrInvalidID             = errors.New("invalid ID")
	ErrParametersMissing     = errors.New("parameters missing")
	ErrInvalidPrice          = errors.New("invalid price")
	ErrInvalidStock          = errors.New("invalid stock")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrProductNotFound       = errors.New("product not found")
	ErrNotFoundProducts      = errors.New("not found products")
	ErrFailedCreatingProduct = errors.New("failed to create product")
	ErrToReduceStock         = errors.New("failed to reduce stock")
	ErrToUpdateProduct       = errors.New("failed to update product")
	ErrToDeletegProduct      = errors.New("failed to delete product")
	ErrScanningRows          = errors.New("failed to scan rows")
)
