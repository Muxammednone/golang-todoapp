package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Погулять с собакой"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("title cant be NULL")
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("title must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("description must be between 1 and 1000")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("completed cant be NULL")
		}
	}

	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask godoc
// @Summary Обновить задачу
// @Description Обновляет информацию об уже существующей в системе задаче
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передвно значение: `"description": "some description"` - устанавливает новое описание в БД
// @Description 3. **Передан null: `"description": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `title` и `completed` не можгут быть выставлены как NULL
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемой задачи"
// @Param request body PatchTaskRequest true "PatchTask тело запроса"
// @Success 200 {object} PatchTaskResponse "Успешно измененная задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
