package main

import (
	"log/slog"
	"os"

	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/config"
	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/lib/logger/sl"
	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) // выводим параметр с адресом
	log.Debug("logger debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // логирование всех запросов
	router.Use(middleware.Recoverer) // если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // парсер URLов поступающих запросов
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
