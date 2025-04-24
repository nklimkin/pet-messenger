package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ru.nklimkin/petmsngr/internal/usecase/user"
)

func NewSignUp(log *slog.Logger, signUp user.UserSignUp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		log.Info("Handle request of user sign up")
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			panic("Invalid id")
		}
		login := chi.URLParam(r, "name")
		signUp.Execute(id, login)
	}
}