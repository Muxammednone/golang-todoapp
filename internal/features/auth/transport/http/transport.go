package auth_transport_http

import (
	"context"

	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/server"
)

type AuthClient interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
}

type AuthHTTPHandler struct {
	authClient     AuthClient
	authMiddleware core_http_middleware.Middleware
}

func NewAuthHTTPHandler(authClient AuthClient, authMiddleware core_http_middleware.Middleware) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authClient:     authClient,
		authMiddleware: authMiddleware,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     "GET",
			Path:       "/auth/me",
			Handler:    h.Me,
			Middleware: []core_http_middleware.Middleware{h.authMiddleware},
		},
		{
			Method:  "POST",
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  "POST",
			Path:    "/auth/register",
			Handler: h.Register,
		},
	}
}
