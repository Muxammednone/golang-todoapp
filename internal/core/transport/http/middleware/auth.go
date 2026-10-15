package core_http_middleware

import (
	"context"
	"net/http"
	"strings"

	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

type TokenVerifier interface {
	Verify(tokenString string) (int64, error)
}

type ssoUserIDKey struct{}

func SSOUserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ssoUserIDKey{}).(int64)
	return id, ok
}

func Auth(verifier TokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := core_logger.FromContext(r.Context())
			response := core_http_response.NewHTTPResponseHandler(logger, w)

			parts := strings.Fields(r.Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.ErrorResponse(
					core_errors.ErrUnauthenticated,
					"authentication required",
				)
				return
			}

			userID, err := verifier.Verify(parts[1])
			if err != nil {
				response.ErrorResponse(
					core_errors.ErrUnauthenticated,
					"invalid or expired token",
				)
				return
			}

			ctx := context.WithValue(r.Context(), ssoUserIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
