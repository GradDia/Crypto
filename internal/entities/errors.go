package entities

import "github.com/pkg/errors"

var (
	ErrInvalidParam = errors.New("invalid param")
	MissingData     = errors.New("missing data")
)
