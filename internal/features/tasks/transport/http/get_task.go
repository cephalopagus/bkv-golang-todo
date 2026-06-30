package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_request "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TaskHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get task id path value",
		)
		return
	}
	taskDomain, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
