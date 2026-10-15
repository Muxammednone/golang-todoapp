package auth_client_grpc

import (
	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return core_errors.ErrInvalidArgument
	case codes.AlreadyExists:
		return core_errors.ErrConflict
	case codes.Unauthenticated:
		return core_errors.ErrUnauthenticated
	case codes.Unavailable:
		return core_errors.ErrUnavailable
	case codes.DeadlineExceeded:
		return core_errors.ErrDeadlineExceeded
	default:
		return err
	}
}
