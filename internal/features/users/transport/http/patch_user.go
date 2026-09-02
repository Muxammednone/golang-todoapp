package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Muxammednone/golang-todoapp/internal/core/domain"
	core_logger "github.com/Muxammednone/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Muxammednone/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Maxim Maximovich"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79876543256"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName cant be NULL")
		}
		fullNameLength := len([]rune(*r.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("FullName len must be between 3 and 100")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLength := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLength < 10 || phoneNumberLength > 15 {
				return fmt.Errorf("PhoneNumber len must be between 10 and 15")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must start with '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передвно значение: `"phone_number": "+79873268732"` - устанавливает новый номер в БД
// @Description 3. **Передан null: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `full_name` не может быть выставлен как NULL
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID изменяемого пользователя"
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Успешно измененный пользователь"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(
		response,
		http.StatusOK,
	)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
