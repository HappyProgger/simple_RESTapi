package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"simple_RESTapi/internal/config"
	logger "simple_RESTapi/internal/http-server"
	"simple_RESTapi/internal/http-server/handlers/save"
	"simple_RESTapi/internal/lib/logger/sl"
	"simple_RESTapi/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Storage struct {
	db *sql.DB
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.App.Env)
	log = log.With(slog.String("env", cfg.App.Env))

	log.Info("initializing server", slog.String("address", cfg.App.Adres)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	storage, err := postgres.New(log)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		log.Info(err.Error())
		panic("Can't connect to the database")
	}
	log.Debug("DB is connected")

	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов

	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	router.Use(logger.New(log))

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("addres", cfg.App.Adres), slog.String("port", strconv.Itoa(cfg.App.Port)))

	fullServerAdr := cfg.App.Adres + ":" + strconv.Itoa(cfg.App.Port)
	srv := &http.Server{
		Addr:    fullServerAdr,
		Handler: router,
		//TODO вынести ReadTimeout WriteTimeout IdleTimeout в конфиг
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error", sl.Err(err))
	}

	log.Error("server stopped")

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
