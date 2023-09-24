package main

import (
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/st-matskevich/audio-guide-bot/api/controller"
)

func main() {
	log.Println("Starting API service")

	app := fiber.New(fiber.Config{
		// Allows to get body as io.Reader/io.Writer
		StreamRequestBody: true,
	})

	// Set logger format to be equal to controller.HandlerPrintf
	app.Use(logger.New(logger.Config{
		Format:        "${time} ${method} ${path}: Returned ${status} in ${latency}\n",
		TimeFormat:    "2006/02/01 15:04:05",
		DisableColors: true,
	}))

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	botInteractor, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatalf("Telegram API initialization error: %v", err)
	}
	log.Println("Telegram API initialized")

	url := os.Getenv("TELEGRAM_WEB_APP_URL")
	controllers := []controller.Controller{
		&controller.BotController{
			BotInteractor: botInteractor,
			WebAppURL:     url,
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
