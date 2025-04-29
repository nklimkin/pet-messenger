package message

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NewMessageRequest struct {
	ChatId string
	Payload string
}

func HandleMessage(log *slog.Logger,) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("can't upgrade connection to websocket, error: %w", err)
			return
		}

		for {

			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Error("can't read message, error: %w", err)
			}

			log.Debug("Handle message: %s with type %s", msg, messageType)

			var newMessage NewMessageRequest
			if err := json.Unmarshal(msg, &newMessage); err != nil {
				log.Error("can't read message, error: %w", err)
				continue
			}

			log.Debug("process message = [%s]", newMessage)
			// нужно будет что то отправить
		}
	}
}
