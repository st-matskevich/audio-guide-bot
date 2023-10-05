package main

import (
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/st-matskevich/audio-guide-bot/api/auth"
	"github.com/st-matskevich/audio-guide-bot/api/blob"
	"github.com/st-matskevich/audio-guide-bot/api/controller"
	"github.com/st-matskevich/audio-guide-bot/api/db"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

func main() {
	log.Println("Starting API service")

	dbURL := os.Getenv("DB_CONNECTION_STRING")
	dbProvider, err := db.CreatePostgresDBProvider(dbURL)
	if err != nil {
		log.Fatalf("PostgreSQL initialization error: %v", err)
	}
	log.Println("PostgreSQL initialized")

	// Apply DB migrations if --migrate is passed
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--migrate" {
		log.Println("Running in DB migration mode")
		err = dbProvider.Migrate()
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		log.Println("Successfully applied DB migrations")
		os.Exit(0)
	}

	app := fiber.New()

	// Setup CORS if CORS_ALLOWED_ORIGINS is provided
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		log.Printf("Setting CORS Allowed Origins to %s", corsOrigins)
		app.Use(cors.New(cors.Config{
			AllowOrigins: corsOrigins,
		}))
	}

	// Set logger format to be equal to controller.HandlerPrintf
	app.Use(logger.New(logger.Config{
		Format:        "${time} ${method} ${path}: Returned ${status} in ${latency}\n",
		TimeFormat:    "2006/02/01 15:04:05",
		DisableColors: true,
	}))

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	botInteractor, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		log.Fatalf("Telegram API initialization error: %v", err)
	}
	log.Println("Telegram API initialized")

	s3URL := os.Getenv("S3_CONNECTION_STRING")
	blobProvider, err := blob.CreateS3BlobProvider(s3URL)
	if err != nil {
		log.Fatalf("S3 blob provider initialization error: %v", err)
	}
	log.Println("S3 blob provider initialized")

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenProvier, err := auth.CreateJWTTokenProvider(jwtSecret)
	if err != nil {
		log.Fatalf("JWT token provider initialization error: %v", err)
	}
	log.Println("JWT token provider initialized")

	webAppURL := os.Getenv("TELEGRAM_WEB_APP_URL")
	paymentsToken := os.Getenv("TELEGRAM_PAYMENTS_TOKEN")
	repository := repository.Repository{DBProvider: dbProvider}
	controllers := []controller.Controller{
		&controller.BotController{
			WebAppURL:        webAppURL,
			BotPaymentsToken: paymentsToken,
			BotInteractor:    botInteractor,
			TicketRepository: &repository,
			ConfigRepository: &repository,
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
	log.Fatal(app.Listen(":" + port))
}
