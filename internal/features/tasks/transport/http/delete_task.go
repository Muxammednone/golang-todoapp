package tasks_transport_http

import (
	"net/http"

	core_errors "github.com/Muxammednone/golang-todoapp/internal/core/errors"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/middleware"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаление существующей задачи в системе по ее ID
// @Tags tasks
// @Param id path int true "ID удаляемой задачи"
// @Success 204 "Успешное удаление задачи"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "Profile missing, task absent or owned by another account"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Failure 401 {object} core_http_response.ErrorResponse "Invalid or expired JWT"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)
	userID, ok := core_http_middleware.ProfileUserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(core_errors.ErrUnauthenticated, "authentication required")
		return
	}

	taskID, err := core_http_request.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)

		return
	}

	responseHandler.NoContentResponse()
}
