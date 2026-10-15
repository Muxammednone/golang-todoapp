package auth_transport_http

import (
	"net/http"

	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
}

// Login godoc
// @Summary Login and receive a JWT
// @Tags auth
// @Produce json
// @Accept json
// @Param request body LoginRequestDTO true "Request body"
// @Success 200 {object} LoginResponseDTO "Success"
// @Failure 400 {object} core_http_response.ErrorResponse "Invalid request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal error"
// @Failure 503 {object} core_http_response.ErrorResponse "SSO unavailable"
// @Failure 504 {object} core_http_response.ErrorResponse "SSO timeout"
// @Failure 401 {object} core_http_response.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	accessToken, err := h.authClient.Login(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to login from authClient",
		)
		return
	}

	responseHandler.JSONResponse(
		LoginResponseDTO{AccessToken: accessToken},
		http.StatusOK,
	)
}
