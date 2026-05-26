package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/cephalopagus/bkv-golang-todo/internal/core/logger"
	core_http_request "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/cephalopagus/bkv-golang-todo/internal/core/transport/http/response"
)

type PatchUserRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func (h *UserHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}
	log.Debug(
		fmt.Sprintf("PatchUserRequest fields:\nFullName:'%s'\nPhoneNumber:'%s'",
			request.FullName,
			request.PhoneNumber),
	)
	w.WriteHeader(http.StatusOK)

}
