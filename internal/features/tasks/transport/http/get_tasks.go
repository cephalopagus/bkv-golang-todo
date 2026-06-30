package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_request "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TaskHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, offset, limit, err := getUserIDLimitOffSetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userId/limit'/'offset' query param",
		)
		return
	}

	taskDomain, err := h.taskService.GetTasks(
		ctx, userID, limit, offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}
func getUserIDLimitOffSetQueryParams(r *http.Request) (*int, *int, *int, error) {

	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limmit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return userID, offset, limit, nil
}
