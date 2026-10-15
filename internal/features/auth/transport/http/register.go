package auth_transport_http

import (
	"net/http"

	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponseDTO struct {
	UserID int64 `json:"user_id"`
}

// Register godoc
// @Summary Register an SSO account
// @Tags auth
// @Produce json
// @Accept json
// @Param request body RegisterRequestDTO true "Request body"
// @Success 201 {object} RegisterResponseDTO "Success"
// @Failure 400 {object} core_http_response.ErrorResponse "Invalid request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal error"
// @Failure 503 {object} core_http_response.ErrorResponse "SSO unavailable"
// @Failure 504 {object} core_http_response.ErrorResponse "SSO timeout"
// @Failure 409 {object} core_http_response.ErrorResponse "Email already registered"
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RegisterRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate",
		)
		return
	}

	userID, err := h.authClient.Register(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"register authClient",
		)
		return
	}

	responseHandler.JSONResponse(
		RegisterResponseDTO{UserID: userID},
		http.StatusCreated,
	)
}
