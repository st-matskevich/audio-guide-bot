package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/st-matskevich/audio-guide-bot/api/bot"
)

type BotController struct {
	BotInteractor bot.BotInteractor
	WebAppURL     string
}

func (controller *BotController) GetRoutes() []Route {
	return []Route{
		{
			Method:  "POST",
			Path:    "/bot",
			Handler: controller.HandleBotUpdate,
		},
	}
}

func (controller *BotController) HandleBotUpdate(c *fiber.Ctx) error {
	update := gotgbot.Update{}
	if err := c.BodyParser(&update); err != nil {
		HandlerPrintf(c, "Failed to parse input - %v", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse input")
	}

	if update.Message == nil {
		HandlerPrintf(c, "Bot update didn't include a message")
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a message")
	}

	message := "Let's start the tour!🎧\nPlease tap the button below to proceed"
	opts := &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Start", WebApp: &gotgbot.WebAppInfo{Url: controller.WebAppURL}},
			}},
		},
	}

	if _, err := controller.BotInteractor.SendMessage(update.Message.Chat.Id, message, opts); err != nil {
		HandlerPrintf(c, "Failed to send bot message - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}
