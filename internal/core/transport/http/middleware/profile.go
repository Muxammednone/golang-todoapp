package core_http_middleware

import (
	"context"
	"net/http"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

type ProfileProvider interface {
	GetProfile(context.Context, int64) (domain.User, error)
}

type profileUserIDKey struct{}

func ProfileUserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(profileUserIDKey{}).(int)
	return id, ok && id > 0
}

// RequireProfile runs after Auth and resolves the SSO account to a local owner.
func RequireProfile(profiles ProfileProvider) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(r.Context()), w)
			id, ok := SSOUserIDFromContext(r.Context())
			if !ok {
				response.ErrorResponse(core_errors.ErrUnauthenticated, "authentication required")
				return
			}
			profile, err := profiles.GetProfile(r.Context(), id)
			if err != nil {
				response.ErrorResponse(err, "create your profile first via POST /api/v1/me")
				return
			}
			ctx := context.WithValue(r.Context(), profileUserIDKey{}, profile.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
