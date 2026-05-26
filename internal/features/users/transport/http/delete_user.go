package users_transport_http

import (
	"net/http"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
	core_http_utils "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeteleUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHadler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHadler.ErrorResponse(err, "failed to get user id path value")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHadler.ErrorResponse(err, "failed to delete user")
		return
	}
	responseHadler.NoContentResponse()
}
