package chat

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"ru.nklimkin/petmsngr/internal/domain/user"
	chat_usecase "ru.nklimkin/petmsngr/internal/usecase/chat"
	"ru.nklimkin/petmsngr/pkg/api/response"
)

type NewChatRequest struct {
	FirstUserId int64 `json:"first_user_id" validate:"required"`
	SecondUserId int64 `json:"second_user_id" validate:"required"`
}

func NewCreateChatHandler(log *slog.Logger, createChat chat_usecase.CreateChat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request NewChatRequest
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error("invalid request body"))
			return
		}

		if err := validator.New().Struct(request); err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error(err.Error()))
			return
		}		
		_, err = createChat.Execute(
			user.UserId{Value: request.FirstUserId},
			user.UserId{Value: request.SecondUserId},
		)
		if err != nil {
			log.Error("error while create chat", slog.Any("error", err))
			response.JSON(r, w, http.StatusInternalServerError, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.Ok())
	}
}

func NewGetUserChatsHandler(log *slog.Logger, getUserChats chat_usecase.GetUserChats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handle request to get chats")
		userId, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
		if err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error("invalid path parameter - user_id"))
			return
		}
		chats, err := getUserChats.Execute(user.UserId{Value: userId})
		if err != nil {
			log.Error("error while get user chats", slog.Any("error", err))
			response.JSON(r, w, http.StatusInternalServerError, response.Error(err.Error()))
			return
		}
		render.JSON(w, r, chats)
	}
}
