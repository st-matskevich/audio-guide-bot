package main

import (
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/st-matskevich/audio-guide-bot/api/auth"
	"github.com/st-matskevich/audio-guide-bot/api/blob"
	"github.com/st-matskevich/audio-guide-bot/api/controller"
	"github.com/st-matskevich/audio-guide-bot/api/db"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

func main() {
	log.Println("Starting API service")

	app := fiber.New()

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

	dbURL := os.Getenv("DB_CONNECTION_STRING")
	dbProvider, err := db.CreatePostgreSQLDBProvider(dbURL)
	if err != nil {
		log.Fatalf("PostgreSQL initialization error: %v", err)
	}
	log.Println("PostgreSQL initialized")

	s3URL := os.Getenv("S3_CONNECTION_STRING")
	blobProvider, err := blob.CreateS3BlobProvider(s3URL)
	if err != nil {
		log.Fatalf("S3 blob provider initialization error: %v", err)
	}
	log.Println("S3 blob provider initialized")

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenProvier := auth.JWTTokenProvider{JWTSecret: []byte(jwtSecret)}

	webAppURL := os.Getenv("TELEGRAM_WEB_APP_URL")
	paymentsToken := os.Getenv("TELEGRAM_PAYMENTS_TOKEN")
	repository := repository.Repository{DBProvider: dbProvider}
	controllers := []controller.Controller{
		&controller.BotController{
			WebAppURL:        webAppURL,
			BotPaymentsToken: paymentsToken,
			BotInteractor:    botInteractor,
			TicketRepository: &repository,
		},
		&controller.TicketsController{
			TokenProvider:    &tokenProvier,
			TicketRepository: &repository,
		},
		&controller.ObjectsController{
			TokenProvider:    &tokenProvier,
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
