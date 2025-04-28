package chat

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
	chat_usecase "ru.nklimkin/petmsngr/internal/usecase/chat"
)

func NewCreateChatHandler(createChat chat_usecase.CreateChat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			panic("Invalid id")
		}
		firstUserId, err := strconv.ParseInt(chi.URLParam(r, "first_user_id"), 10, 64)
		if err != nil {
			panic("Invalid first_user_id")
		}
		secondUserId, err := strconv.ParseInt(chi.URLParam(r, "second_user_id"), 10, 64)
		if err != nil {
			panic("Invalid first_user_id")
		}
		createChat.Execute(
			chat.ChatId{Value: chatId},
			user.UserId{Value: firstUserId},
			user.UserId{Value: secondUserId},
		)
	}
}

func NewGetUserChatsHandler(log *slog.Logger, getUserChats chat_usecase.GetUserChats) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Handle request to get chats")
		userId, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
		if err != nil {
			
		}
		chats := getUserChats.Execute(user.UserId{Value: userId})
		render.JSON(w, r, chats)
	}
}
