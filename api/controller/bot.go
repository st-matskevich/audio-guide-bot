package controller

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/bot"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

const BUY_TICKET_QUERY = "buy_ticket"
const TICKET_CURRENCY = "USD"
const TICKET_PRICE = 100

type BotController struct {
	WebAppURL        string
	BotPaymentsToken string
	BotInteractor    bot.BotInteractor
	TicketRepository repository.TicketRepository
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

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotMessage(c *fiber.Ctx, update *gotgbot.Update) error {
	if update.Message.SuccessfulPayment != nil {
		ticketCode, err := uuid.Parse(update.Message.SuccessfulPayment.InvoicePayload)
		if err != nil {
			HandlerPrintf(c, "Failed to parse bot payment payload - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to parse bot payment payload")
		}

		if err = controller.TicketRepository.CreateTicket(ticketCode.String()); err != nil {
			HandlerPrintf(c, "Failed to register ticket in DB - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to register ticket in DB")
		}

		message := controller.buildPurchaseMessage(ticketCode)
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
		ticketCode := uuid.New()
		if _, err := controller.BotInteractor.SendInvoice(update.CallbackQuery.Message.Chat.Id, title, description, ticketCode.String(), controller.BotPaymentsToken, TICKET_CURRENCY, price, nil); err != nil {
			HandlerPrintf(c, "Failed to send bot invoice - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot invoice")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

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
	message := "Let's start the tour!🎧\nPlease choose an option below to proceed"
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

func (controller *BotController) buildPurchaseMessage(ticketCode uuid.UUID) BotMessage {
	ticketString := ticketCode.String()
	message := "Thank you for your purchase!\nYour ticket number: " + ticketString + "\nPlease tap the button below to proceed"
	appURL := controller.WebAppURL + "?ticket=" + ticketString
	opts := &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Start the tour", WebApp: &gotgbot.WebAppInfo{Url: appURL}},
			}},
		},
	}

	return BotMessage{Message: message, Options: opts}
}
