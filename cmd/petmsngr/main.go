package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ru.nklimkin/petmsngr/internal/config"
	user_in_memmory_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/user"
	user_postgres_persistence "ru.nklimkin/petmsngr/internal/peristence/postgres/user"
	"ru.nklimkin/petmsngr/internal/router/user"
	user_usecase "ru.nklimkin/petmsngr/internal/usecase/user"

	chat_in_memmory_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/chat"
	chat_postgres_persistence "ru.nklimkin/petmsngr/internal/peristence/postgres/chat"
	"ru.nklimkin/petmsngr/internal/router/chat"
	chat_usecase "ru.nklimkin/petmsngr/internal/usecase/chat"

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
	POSTGRES = "postgres"
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

	setupUserHandlers(log, cfg, router)
	setupChatHandlers(log, cfg, router)
	setupMessageHandlers(router)

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

func setupUserHandlers(log *slog.Logger, cfg *config.Config, router *chi.Mux) {
	// userAccessor := buildUserAccessor(cfg.Datasource)
	userPersistence := buildUserPersistence(cfg.Datasource)
	userSignUpUsecase := user_usecase.New(userPersistence)

	router.Post("/api/v1/user/sign-up", user.NewSignUp(log, userSignUpUsecase))
}

func setupChatHandlers(log *slog.Logger, cfg *config.Config, router *chi.Mux) {
	chatAccessor := buildChatAccessor(cfg.Datasource)
	chatPersistence := buildChatPersistence(cfg.Datasource)

	getChatUsercase := chat_usecase.NewGetUserChats(chatAccessor)
	createChatUsecase := chat_usecase.NewCreateChat(chatPersistence)

	router.Post("/api/v1/chat", chat.NewCreateChatHandler(log, createChatUsecase))
	router.Get("/api/v1/chat/user/{user_id}", chat.NewGetUserChatsHandler(log, getChatUsercase))
}

func setupMessageHandlers(router *chi.Mux) {
	// router.HandleFunc("/ws", message.HandleMessage())
}

func buildChatAccessor(datasource config.Datasource) chat_usecase.ChatAccessor {

	switch datasource.Type {
	case POSTGRES:
		rep, err := chat_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic("Can't build chat postgres repository")
		}
		return rep
	case IN_MEMMORY:
		return chat_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}

func buildChatPersistence(datasource config.Datasource) chat_usecase.ChatPersistence {

	switch datasource.Type {
	case POSTGRES:
		rep, err := chat_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic("Can't build chat postgres repository")
		}
		return rep
	case IN_MEMMORY:
		return chat_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}

func buildUserAccessor(datasource config.Datasource) user_usecase.UserAccessor {

	switch datasource.Type {
	case POSTGRES:
		rep, err := user_postgres_persistence.New(datasource)
		if err != nil {
			panic("Can't build user postgres repository")
		}
		return rep
	case IN_MEMMORY:
		return user_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}

func buildUserPersistence(datasource config.Datasource) user_usecase.UserPersistence {

	switch datasource.Type {
	case POSTGRES:
		rep, err := user_postgres_persistence.New(datasource)
		if err != nil {
			panic(err)
		}
		return rep
	case IN_MEMMORY:
		return user_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}

func buildMessageAccessor(datasource config.Datasource) message_usecase.ChatMessageAccessor {

	switch datasource.Type {
	case POSTGRES:
		rep, err := message_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic("Can't build message postgres repository")
		}
		return rep
	case IN_MEMMORY:
		return message_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}

func buildMessagePersistence(datasource config.Datasource) message_usecase.ChatMessagePersistence {

	switch datasource.Type {
	case POSTGRES:
		rep, err := message_postgres_persistence.NewPostgresRepository(datasource)
		if err != nil {
			panic("Can't build message postgres repository")
		}
		return rep
	case IN_MEMMORY:
		return message_in_memmory_persistence.New()
	}
	panic("Invalid datasource type")
}
