package auth_transport_http

import (
	"net/http"

	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

type MeResponseDTO struct {
	SSOUserID int64 `json:"sso_user_id"`
}

// Me godoc
// @Summary Get SSO account ID from JWT; no profile required
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponseDTO "Success"
// @Failure 401 {object} core_http_response.ErrorResponse "Invalid or expired JWT"
// @Router /auth/me [get]
func (h *AuthHTTPHandler) Me(w http.ResponseWriter, r *http.Request) {
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(r.Context()), w)
	userID, ok := core_http_middleware.SSOUserIDFromContext(r.Context())
	if !ok {
		response.ErrorResponse(core_errors.ErrUnauthenticated, "authentication required")
		return
	}
	response.JSONResponse(MeResponseDTO{SSOUserID: userID}, http.StatusOK)
}
