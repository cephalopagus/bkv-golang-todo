package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/cephalopagus/bkv-golang-todo/internal/core/domain"
	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_request "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
	core_http_types "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'Title cannot be NULL'")
		}
		titleLenght := len([]rune(*r.Title.Value))
		if titleLenght < 0 || titleLenght > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLenght := len([]rune(*r.Description.Value))
			if descriptionLenght < 1 || descriptionLenght > 1000 {
				return fmt.Errorf("'Descriprion' must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("'Completed' cannot be NULL")
		}
	}
	return nil
}

type PatchUserResponse TaskDTOResponse

func (h *TaskHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get id path value",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.taskService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}
	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
