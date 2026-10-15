package core_errors

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrConflict         = errors.New("conflict")
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrUnavailable      = errors.New("unavailable")
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)
