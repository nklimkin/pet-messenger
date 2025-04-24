package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ru.nklimkin/petmsngr/internal/config"
	user_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/user"
	"ru.nklimkin/petmsngr/internal/router/user"
	user_usecase "ru.nklimkin/petmsngr/internal/usecase/user"

	chat_persistence "ru.nklimkin/petmsngr/internal/peristence/in_memmory/chat"
	"ru.nklimkin/petmsngr/internal/router/chat"
	chat_usecase "ru.nklimkin/petmsngr/internal/usecase/chat"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.Load()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Startup web service")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	setupUserHandlers(log, router)
	setupChatHandlers(log, router)
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

func setupUserHandlers(log *slog.Logger, router *chi.Mux) {
	userRepository := user_persistence.New()
	userSignUpUsecase := user_usecase.New(userRepository)

	router.Post("/api/v1/user/sign-up", user.NewSignUp(log, userSignUpUsecase))
}

func setupChatHandlers(log *slog.Logger, router *chi.Mux) {
	chatRepository := chat_persistence.New()

	getChatUsercase := chat_usecase.NewGetUserChats(chatRepository)
	createChatUsecase := chat_usecase.NewCreateChat(chatRepository)

	router.Post("/api/v1/chat", chat.NewCreateChatHandler(createChatUsecase))
	router.Get("/api/v1/chat/user/{user_id}", chat.NewGetUserChatsHandler(log, getChatUsercase))
}

func setupMessageHandlers(router *chi.Mux) {
	// router.HandleFunc("/ws", message.HandleMessage())
}
