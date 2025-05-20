package entities

import "github.com/pkg/errors"

var (
	ErrInvalidParam = errors.New("invalid param")
	ErrInternal     = errors.New("internal error")
	ErrNotFound     = errors.New("missing data")
	ErrInvalidInput = errors.New("invalid input")
)
