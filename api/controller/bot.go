package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/bot"
)

const BUY_TICKET_QUERY = "buy_ticket"
const TICKET_CURRENCY = "USD"
const TICKET_PRICE = 100

type BotController struct {
	BotInteractor    bot.BotInteractor
	WebAppURL        string
	BotPaymentsToken string
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

	if update.CallbackQuery != nil {
		return controller.HandleBotCallback(c, &update)
	}

	if update.PreCheckoutQuery != nil {
		return controller.HandleBotPreCheckout(c, &update)
	}

	if update.Message != nil {
		return controller.HandleBotMessage(c, &update)
	}

	// TODO: should we return error in case if no handlers were called?
	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotMessage(c *fiber.Ctx, update *gotgbot.Update) error {
	if update.Message.SuccessfulPayment != nil {
		ticketID, err := uuid.Parse(update.Message.SuccessfulPayment.InvoicePayload)
		if err != nil {
			HandlerPrintf(c, "Failed to parse bot payment payload - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to parse bot payment payload")
		}

		message := controller.buildPurchaseMessage(ticketID)
		if _, err := controller.BotInteractor.SendMessage(update.Message.Chat.Id, message.Message, message.Options); err != nil {
			HandlerPrintf(c, "Failed to send bot message - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	message := controller.buildWelcomeMessage()
	if _, err := controller.BotInteractor.SendMessage(update.Message.Chat.Id, message.Message, message.Options); err != nil {
		HandlerPrintf(c, "Failed to send bot message - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotCallback(c *fiber.Ctx, update *gotgbot.Update) error {
	if update.CallbackQuery.Data == "" {
		HandlerPrintf(c, "Bot update didn't include a callback data")
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback data")
	}

	if result, err := controller.BotInteractor.AnswerCallbackQuery(update.CallbackQuery.Id, nil); !result || err != nil {
		HandlerPrintf(c, "Failed to answer callback query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer callback query")
	}

	if update.CallbackQuery.Data == BUY_TICKET_QUERY {
		if update.CallbackQuery.Message == nil {
			HandlerPrintf(c, "Bot update didn't include a callback message")
			return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback message")
		}

		title := "Tour ticket"
		description := "Ticket that allows to start the tour"
		price := []gotgbot.LabeledPrice{{Label: "Price", Amount: int64(TICKET_PRICE)}}
		ticketID := uuid.New()
		if _, err := controller.BotInteractor.SendInvoice(update.CallbackQuery.Message.Chat.Id, title, description, ticketID.String(), controller.BotPaymentsToken, TICKET_CURRENCY, price, nil); err != nil {
			HandlerPrintf(c, "Failed to send bot invoice - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot invoice")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	// TODO: should we return error in case if no handlers were called?
	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotPreCheckout(c *fiber.Ctx, update *gotgbot.Update) error {
	acceptCheckout := true

	if update.PreCheckoutQuery.Currency != TICKET_CURRENCY {
		acceptCheckout = false
	}

	if update.PreCheckoutQuery.TotalAmount != TICKET_PRICE {
		acceptCheckout = false
	}

	if _, err := uuid.Parse(update.PreCheckoutQuery.InvoicePayload); err != nil {
		acceptCheckout = false
	}

	if result, err := controller.BotInteractor.AnswerPreCheckoutQuery(update.PreCheckoutQuery.Id, acceptCheckout, nil); !result || err != nil {
		HandlerPrintf(c, "Failed to answer pre-checkout query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer pre-checkout query")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

type BotMessage struct {
	Message string
	Options *gotgbot.SendMessageOpts
}

func (controller *BotController) buildWelcomeMessage() BotMessage {
	message := "Let's start the tour!🎧\nPlease tap the button below to proceed"
	opts := &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Start the tour", WebApp: &gotgbot.WebAppInfo{Url: controller.WebAppURL}},
			}, {
				{Text: "Buy a ticket", CallbackData: BUY_TICKET_QUERY},
			}},
		},
	}

	return BotMessage{Message: message, Options: opts}
}

func (controller *BotController) buildPurchaseMessage(ticketID uuid.UUID) BotMessage {
	message := "Thank you for your purchase!\nYour ticket number: " + ticketID.String() + "\nPlease tap the button below to proceed"
	opts := &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Start the tour", WebApp: &gotgbot.WebAppInfo{Url: controller.WebAppURL}},
			}},
		},
	}

	return BotMessage{Message: message, Options: opts}
}
