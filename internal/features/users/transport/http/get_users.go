package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
	core_http_utils "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UserHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHadler := core_http_response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffSetQueryParams(r)
	if err != nil {
		responseHadler.ErrorResponse(err, "failed to get 'limit'/'offset' query param")
		return
	}
	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHadler.ErrorResponse(err, "failed to get users")
		return
	}
	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHadler.JSONResponse(response, http.StatusOK)

}

func getLimitOffSetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limmit' query param: %w", err)
	}
	offset, err := core_http_utils.GetIntQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return limit, offset, nil
}
