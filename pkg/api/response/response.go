package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func Error(msg string) Response {
	return Response{StatusError, msg}
}

func Ok() Response {
	return Response{Status: StatusOk}
}

func JSON(
	r *http.Request,
	w http.ResponseWriter,
	status int,
	payload interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	render.JSON(w, r, payload)
}
