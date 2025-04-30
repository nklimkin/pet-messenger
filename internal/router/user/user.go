package user

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"ru.nklimkin/petmsngr/internal/usecase/user"
	"ru.nklimkin/petmsngr/pkg/api/response"
)

type SignUpRequest struct {
	Id    int64  `json:"id" validate:"required"`
	Login string `json:"login" validate:"required"`
}

func NewSignUp(log *slog.Logger, signUp user.UserSignUp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handle request of user sign up")
		var signUpRequest SignUpRequest
		err := render.DecodeJSON(r.Body, &signUpRequest)
		if err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error("invalid request body"))
			return
		}

		if err := validator.New().Struct(signUpRequest); err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error(err.Error()))
			return
		}

		_, err = signUp.Execute(signUpRequest.Id, signUpRequest.Login)
		if err != nil {
			log.Error("Failed to sigup user", slog.Any("error", err))
			response.JSON(r, w, http.StatusInternalServerError, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.Ok())
	}
}
