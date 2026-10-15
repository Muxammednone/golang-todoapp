package users_transport_http

import (
	"net/http"

	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

// CreateProfile godoc
// @Summary Create your TodoApp profile
// @Tags profile
// @Produce json
// @Accept json
// @Param request body CreateUserRequest true "Request body"
// @Security BearerAuth
// @Success 201 {object} UserDTOResponse "Success"
// @Failure 400 {object} core_http_response.ErrorResponse "Invalid profile"
// @Failure 401 {object} core_http_response.ErrorResponse "JWT required"
// @Failure 409 {object} core_http_response.ErrorResponse "Profile already exists"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal error"
// @Router /me [post]
func (h *UsersHTTPHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), w)
	ssoUserID, ok := core_http_middleware.SSOUserIDFromContext(ctx)
	if !ok {
		response.ErrorResponse(core_errors.ErrUnauthenticated, "authentication required")
		return
	}
	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		response.ErrorResponse(err, "invalid profile")
		return
	}
	user, err := h.usersService.CreateProfile(ctx, ssoUserID, domainFromDTO(request))
	if err != nil {
		response.ErrorResponse(err, "failed to create profile")
		return
	}
	response.JSONResponse(userDTOFromDomain(user), http.StatusCreated)
}

// GetProfile godoc
// @Summary Get your TodoApp profile
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserDTOResponse "Success"
// @Failure 401 {object} core_http_response.ErrorResponse "JWT required"
// @Failure 404 {object} core_http_response.ErrorResponse "Create a profile first"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal error"
// @Router /me [get]
func (h *UsersHTTPHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), w)
	ssoUserID, ok := core_http_middleware.SSOUserIDFromContext(ctx)
	if !ok {
		response.ErrorResponse(core_errors.ErrUnauthenticated, "authentication required")
		return
	}
	user, err := h.usersService.GetProfile(ctx, ssoUserID)
	if err != nil {
		response.ErrorResponse(err, "failed to get profile")
		return
	}
	response.JSONResponse(userDTOFromDomain(user), http.StatusOK)
}
