package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ru.nklimkin/petmsngr/internal/config"
	"ru.nklimkin/petmsngr/internal/domain/chat"
	user_in_memmory_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/user"
	user_postgres_persistence "ru.nklimkin/petmsngr/internal/peristence/postgres/user"
	"ru.nklimkin/petmsngr/internal/router/user"
	user_usecase "ru.nklimkin/petmsngr/internal/usecase/user"

	chat_in_memmory_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/chat"
	chat_postgres_persistence "ru.nklimkin/petmsngr/internal/peristence/postgres/chat"
	chat_router "ru.nklimkin/petmsngr/internal/router/chat"
	chat_usecase "ru.nklimkin/petmsngr/internal/usecase/chat"

	message_router "ru.nklimkin/petmsngr/internal/router/message"
	message_in_memmory_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/message"
	message_postgres_persistence "ru.nklimkin/petmsngr/internal/peristence/postgres/message"

	// "ru.nklimkin/petmsngr/internal/router/message"
	message_usecase "ru.nklimkin/petmsngr/internal/usecase/message"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const (
	POSTGRES   = "postgres"
	IN_MEMMORY = "in_memmory"
)

func main() {

	cfg := config.Load()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Startup web service")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	userRequestComponents := buildUserRequestsComponents(cfg.Datasource)
	chatRequestComponents := buildChatRequestsComponents(cfg.Datasource)

	setupUserHandlers(log, router, userRequestComponents)
	setupChatHandlers(log, router, chatRequestComponents, userRequestComponents)
	setupMessageHandlers(log, router)

	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		log.Error("Error while start up server")
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupUserHandlers(
	log *slog.Logger,
	router *chi.Mux,
	userRequestComponents HandleUserRequestComponents,
) {
	userSignUpUsecase := user_usecase.New(userRequestComponents.persistence, userRequestComponents.idGenerator)

	router.Post("/api/v1/user/sign-up", user.NewSignUp(log, userSignUpUsecase))
}

func setupChatHandlers(
	log *slog.Logger,
	router *chi.Mux,
	chatRequestComponents HandleChatRequestComponents,
	userRequestComponents HandleUserRequestComponents,
) {

	getChatUsercase := chat_usecase.NewGetUserChats(chatRequestComponents.accessor)
	createChatUsecase := chat_usecase.NewCreateChat(
		chatRequestComponents.persistence,
		chatRequestComponents.idGenerator,
		userRequestComponents.accessor)

	router.Post("/api/v1/chat", chat_router.NewCreateChatHandler(log, createChatUsecase))
	router.Get("/api/v1/chat/user/{user_id}", chat_router.NewGetUserChatsHandler(log, getChatUsercase))
}

func setupMessageHandlers(log *slog.Logger, router *chi.Mux) {
	clientStorage := message_router.ClientStorage{
		Register: make(chan *message_router.Client),
		Clients: make(map[chat.ChatId][]*message_router.Client, 0),
		Broadcast: make(chan message_router.NewMessageRequest),
	}
	go clientStorage.Init()
	router.HandleFunc("/ws", message_router.HandleMessage(log, &clientStorage))
}

type HandleUserRequestComponents struct {
	accessor    user_usecase.UserAccessor
	persistence user_usecase.UserPersistence
	idGenerator user_usecase.UserIdGenerator
}

type HandleMessageRequestComponents struct {
	accessor    message_usecase.ChatMessageAccessor
	persistence message_usecase.ChatMessagePersistence
	idGenerator message_usecase.MessageIdGenerator
}

type HandleChatRequestComponents struct {
	accessor    chat_usecase.ChatAccessor
	persistence chat_usecase.ChatPersistence
	idGenerator chat_usecase.ChatIdGenerator
}

func buildUserRequestsComponents(datasource config.Datasource) HandleUserRequestComponents {
	switch datasource.Type {
	case POSTGRES:
		rep, err := user_postgres_persistence.New(datasource)
		if err != nil {
			panic(err)
		}
		return HandleUserRequestComponents{rep, rep, rep}
	case IN_MEMMORY:
		rep := user_in_memmory_persistence.New()
		return HandleUserRequestComponents{rep, rep, rep}
	}
	panic("Invalid datasource type")
}

func buildChatRequestsComponents(datasource config.Datasource) HandleChatRequestComponents {
	switch datasource.Type {
	case POSTGRES:
		rep, err := chat_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic(err)
		}
		return HandleChatRequestComponents{rep, rep, rep}
	case IN_MEMMORY:
		rep := chat_in_memmory_persistence.New()
		return HandleChatRequestComponents{rep, rep, rep}
	}
	panic("Invalid datasource type")
}

func buildMessageRequestsComponents(datasource config.Datasource) HandleMessageRequestComponents {
	switch datasource.Type {
	case POSTGRES:
		rep, err := message_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic(err)
		}
		return HandleMessageRequestComponents{rep, rep, rep}
	case IN_MEMMORY:
		rep := message_in_memmory_persistence.New()
		return HandleMessageRequestComponents{rep, rep, rep}
	}
	panic("Invalid datasource type")
}
