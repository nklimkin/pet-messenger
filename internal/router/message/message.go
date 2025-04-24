package message

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageDto struct {
	ChatId string
	Payload string
}

func HandleMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		for {

			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("Error while read message")
			}

			log.Printf("Handle message: %s with type %s", msg, messageType)

			var messageDto MessageDto
			if err := json.Unmarshal(msg, &messageDto); err != nil {
				log.Fatal("Invalid request")
				continue
			}

			log.Print(messageDto)
			// нужно будет что то отправить
		}
	}
}
