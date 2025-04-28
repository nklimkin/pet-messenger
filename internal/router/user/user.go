package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"ru.nklimkin/petmsngr/internal/usecase/user"
)

func NewSignUp(log *slog.Logger, signUp user.UserSignUp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handle request of user sign up")
		query := r.URL.Query()
		id, err := strconv.ParseInt(query.Get("id"), 10, 64)
		if err != nil {
			panic(err)
		}
		login := query.Get("login")
		_, err = signUp.Execute(id, login)
		if err != nil {
			log.Error("Failed to sigup user: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error"))
		}
	}
}
