package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService   UsersService
	authMiddleware core_http_middleware.Middleware
}

type UsersService interface {
	CreateProfile(context.Context, int64, domain.User) (domain.User, error)
	GetProfile(context.Context, int64) (domain.User, error)
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

func NewUsersHTTPHandler(
	usersService UsersService,
	authMiddleware core_http_middleware.Middleware,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService:   usersService,
		authMiddleware: authMiddleware,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/me",
			Handler:    h.CreateProfile,
			Middleware: []core_http_middleware.Middleware{h.authMiddleware},
		},
		{
			Method:     http.MethodGet,
			Path:       "/me",
			Handler:    h.GetProfile,
			Middleware: []core_http_middleware.Middleware{h.authMiddleware},
		},
	}
}
