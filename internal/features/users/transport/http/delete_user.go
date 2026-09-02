package users_transport_http

import (
	"net/http"

	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление существующего в системе пользователя по его ID
// @Tags users
// @Param id path int true "ID удаляемого пользователя"
// @Success 204 "Успешное удаление пользоваетля"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()

}
