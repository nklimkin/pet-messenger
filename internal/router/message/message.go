package message

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	"ru.nklimkin/petmsngr/internal/domain/user"
	"ru.nklimkin/petmsngr/pkg/api/response"
)

type Client struct {
	Connection *websocket.Conn
	ChatId     chat.ChatId
	UserId     user.UserId
	Send       chan string
	Storage    *ClientStorage
}

func (c *Client) readMessage(log *slog.Logger) {
	defer c.Connection.Close()

	for {
		_, msg, err := c.Connection.ReadMessage()
		if err != nil {
			log.Error("can't read request message", slog.Any("error", err))
			break
		}

		var newMessage NewMessageRequest
		if err := json.Unmarshal(msg, &newMessage); err != nil {
			log.Error("can't process request message", slog.Any("error", err))
			break
		}

		c.Storage.Broadcast <- newMessage
	}

}

func (c *Client) writeMessage(log *slog.Logger) {
	defer c.Connection.Close()

	for message := range c.Send {
		if err := c.Connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Error("can't write message", slog.Any("error", err))
			break
		}
	}
}

type ClientStorage struct {
	Register  chan *Client
	Clients   map[chat.ChatId][]*Client
	Broadcast chan NewMessageRequest
}

type NewMessageRequest struct {
	UserId  int64  `json:"user_id"`
	ChatId  int64  `json:"chat_id"`
	Payload string `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (cs *ClientStorage) Init() {
	for {
		select {
		case client := <-cs.Register:
			if client != nil {
				clientsInChat := cs.Clients[client.ChatId]
				if clientsInChat == nil {
					clientsInChat = make([]*Client, 0)
				}
				clientsInChat = append(clientsInChat, client)
				cs.Clients[client.ChatId] = clientsInChat
			}

		case message := <-cs.Broadcast:
			chatId := chat.ChatId{Value: message.ChatId}
			currentUserId := user.UserId{Value: message.UserId}
			chatClients := cs.Clients[chatId]
			for _, chatClient := range chatClients {
				if chatClient.UserId != currentUserId {
					chatClient.Send <- message.Payload
				}
			}
		}
	}
}

func HandleMessage(log *slog.Logger, clientStorage *ClientStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.ParseInt(r.URL.Query().Get("chat_id"), 10, 64)
		if err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error("invalid query parameter - chat_id"))
			return
		}
		userId, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
		if err != nil {
			response.JSON(r, w, http.StatusBadRequest, response.Error("invalid query parameter - user_id"))
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("can't upgrade connection to websocket", slog.Any("error", err))
			return
		}

		client := &Client{
			Connection: conn,
			ChatId:     chat.ChatId{Value: chatId},
			UserId:     user.UserId{Value: userId},
			Send:       make(chan string),
			Storage:    clientStorage,
		}

		clientStorage.Register <- client

		go client.readMessage(log)
		go client.writeMessage(log)
	}
}
