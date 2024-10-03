package main

import (
	//_ "contact-api/docs" // Импортируем сгенерированные документы Swagger
	"contact-api/internal/app/config"
	deleteAll "contact-api/internal/app/http-server/handlers/all/delete"
	getAll "contact-api/internal/app/http-server/handlers/all/get"
	"contact-api/internal/app/http-server/handlers/all/save"
	deleteOne "contact-api/internal/app/http-server/handlers/one/delete"
	getOne "contact-api/internal/app/http-server/handlers/one/get"
	"contact-api/internal/app/http-server/handlers/one/update"
	"contact-api/internal/app/storage/mongo"
	"contact-api/internal/pkg/logger/handlers/slogpretty"
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

// @title My Swagger API
// @version 1.0
// @description Swagger API for Golang Project.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @BasePath /v1/contact

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

func main() {
	cfg := config.MustLoad("./config/prod.yaml")

	log := SetupLogger(cfg.Env)

	log.Info("Starting server at port", slog.String("port", cfg.Port))

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Применяем CORS middleware
	router.Use(c.Handler)

	ctx := context.Background()

	storage, err := mongo.New(log, ctx, cfg.DBConnection)
	if err != nil {
		log.Error("error connecting to database", err)
		panic(err)
	}
	defer storage.Close()

	// Настройка Swagger
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/v1/contact", func(r chi.Router) {
		r.Get("/", getAll.New(log, storage))
		r.Post("/", save.New(log, storage))
		r.Delete("/", deleteAll.New(log, storage))

		r.Route("/{uid}", func(r chi.Router) {
			r.Get("/", getOne.New(log, storage))
			r.Delete("/", deleteOne.New(log, storage))
			r.Put("/", update.New(log, storage))
		})
	})

	err = http.ListenAndServe(cfg.Port, router)
	if err != nil {
		log.Error("Error starting server", err)
	}
}

func SetupLogger(env string) *slog.Logger {
	log := &slog.Logger{}

	switch env {
	case EnvLocal:
		log = setupPrettySlog()
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // Если конфигурация окружения недействительна, по умолчанию устанавливаются параметры prod из соображений безопасности
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
