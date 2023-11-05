package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/st-matskevich/audio-guide-bot/api/controller"
	"github.com/st-matskevich/audio-guide-bot/api/provider/auth"
	"github.com/st-matskevich/audio-guide-bot/api/provider/blob"
	"github.com/st-matskevich/audio-guide-bot/api/provider/bot"
	"github.com/st-matskevich/audio-guide-bot/api/provider/db"
	"github.com/st-matskevich/audio-guide-bot/api/provider/translation"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting API service")

	dbURL := os.Getenv("DB_CONNECTION_STRING")
	dbProvider, err := db.CreatePostgresDBProvider(dbURL)
	if err != nil {
		slog.Error("PostgreSQL initialization error", "error", err)
		os.Exit(1)
	}
	slog.Info("PostgreSQL initialized")

	// Apply DB migrations if --migrate is passed
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--migrate" {
		slog.Info("Running in DB migration mode")
		err = dbProvider.Migrate()
		if err != nil {
			slog.Error("Migration failed", "error", err)
			os.Exit(1)
		}

		slog.Info("Successfully applied DB migrations")
		os.Exit(0)
	}

	app := fiber.New()

	// Setup CORS if CORS_ALLOWED_ORIGINS is provided
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		slog.Info("Setting CORS Allowed Origins", "origins", corsOrigins)
		app.Use(cors.New(cors.Config{
			AllowOrigins: corsOrigins,
		}))
	}

	// Set logger format to be equal to controller.HandlerPrintf
	app.Use(slogfiber.New(slog.Default()))
	/*app.Use(logger.New(logger.Config{
		Format:        "${time} ${method} ${path}: Returned ${status} in ${latency}\n",
		TimeFormat:    "2006/02/01 15:04:05",
		DisableColors: true,
	}))*/

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	paymentsToken := os.Getenv("TELEGRAM_PAYMENTS_TOKEN")
	botProvider, err := bot.CreateTelegramBotProvider(botToken, paymentsToken)
	if err != nil {
		slog.Error("Telegram API initialization error", "error", err)
		os.Exit(1)
	}
	slog.Info("Telegram API initialized")

	s3URL := os.Getenv("S3_CONNECTION_STRING")
	blobProvider, err := blob.CreateS3BlobProvider(s3URL)
	if err != nil {
		slog.Error("S3 blob provider initialization error", "error", err)
		os.Exit(1)
	}
	slog.Info("S3 blob provider initialized")

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenProvier, err := auth.CreateJWTTokenProvider(jwtSecret)
	if err != nil {
		slog.Error("JWT token provider initialization error", "error", err)
		os.Exit(1)
	}
	slog.Info("JWT token provider initialized")

	translationProvier, err := translation.CreateI18NTranslationProvider()
	if err != nil {
		slog.Error("Translation provider initialization error", "error", err)
		os.Exit(1)
	}
	slog.Info("Translation provider initialized")

	webAppURL := os.Getenv("TELEGRAM_WEB_APP_URL")
	repository := repository.Repository{DBProvider: dbProvider}
	controllers := []controller.Controller{
		&controller.BotController{
			WebAppURL:           webAppURL,
			BotProvider:         botProvider,
			TranslationProvider: translationProvier,
			TicketRepository:    &repository,
			ConfigRepository:    &repository,
		},
		&controller.TicketsController{
			TokenProvider:    tokenProvier,
			TicketRepository: &repository,
		},
		&controller.ObjectsController{
			TokenProvider:    tokenProvier,
			BlobProvider:     blobProvider,
			ObjectRepository: &repository,
		},
	}

	for _, controller := range controllers {
		for _, route := range controller.GetRoutes() {
			app.Add(route.Method, route.Path, route.Handler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err = app.Listen(":" + port)
	slog.Error("API HTTP server exited", "error", err)
}
