package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/bot"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

const BUY_TICKET_QUERY = "buy_ticket"
const TICKET_CURRENCY_KEY = "TICKET_CURRENCY"
const TICKET_PRICE_KEY = "TICKET_PRICE"

type BotController struct {
	WebAppURL        string
	BotProvider      bot.BotProvider
	TicketRepository repository.TicketRepository
	ConfigRepository repository.ConfigRepository
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
	update := bot.Update{}
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

func (controller *BotController) HandleBotMessage(c *fiber.Ctx, update *bot.Update) error {
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

		message, options := controller.buildPurchaseMessage(ticketCode)
		if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
			HandlerPrintf(c, "Failed to send bot message - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	message, options := controller.buildWelcomeMessage()
	if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
		HandlerPrintf(c, "Failed to send bot message - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotCallback(c *fiber.Ctx, update *bot.Update) error {
	if update.CallbackQuery.Data == "" {
		HandlerPrintf(c, "Bot update didn't include a callback data")
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback data")
	}

	if err := controller.BotProvider.AnswerCallbackQuery(update.CallbackQuery.Id); err != nil {
		HandlerPrintf(c, "Failed to answer callback query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer callback query")
	}

	if update.CallbackQuery.Data == BUY_TICKET_QUERY {
		if update.CallbackQuery.Message == nil {
			HandlerPrintf(c, "Bot update didn't include a callback message")
			return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback message")
		}

		price, err := controller.getTicketPrice(c)
		if err != nil {
			HandlerPrintf(c, "Failed to get ticket price - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get ticket price")
		} else if price == nil {
			HandlerPrintf(c, "Ticket price is not set")

			message, options := controller.buildPaymentsDisabledMessage()
			if err := controller.BotProvider.SendMessage(update.CallbackQuery.Message.Chat.Id, message, options); err != nil {
				HandlerPrintf(c, "Failed to send bot message - %v", err)
				return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
			}

			return HandlerSendSuccess(c, fiber.StatusOK, nil)
		}

		title := "Tour ticket"
		description := "Ticket that allows to start the tour"
		labeledPrice := bot.InvoicePrice{Currency: price.Currency, Parts: []bot.PricePart{{Label: "Price", Amount: price.Price}}}
		ticketCode := uuid.New()
		if err := controller.BotProvider.SendInvoice(update.CallbackQuery.Message.Chat.Id, title, description, ticketCode.String(), labeledPrice); err != nil {
			HandlerPrintf(c, "Failed to send bot invoice - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot invoice")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotPreCheckout(c *fiber.Ctx, update *bot.Update) error {
	acceptCheckout, errorMessage, err := controller.validatePreCheckoutQuery(c, update)
	if err != nil {
		HandlerPrintf(c, "Failed to validate pre-checkout query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to validate pre-checkout query")
	}

	if err := controller.BotProvider.AnswerPreCheckoutQuery(update.PreCheckoutQuery.Id, acceptCheckout, bot.AnswerPreCheckoutQueryOptions{ErrorMessage: errorMessage}); err != nil {
		HandlerPrintf(c, "Failed to answer pre-checkout query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer pre-checkout query")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) validatePreCheckoutQuery(c *fiber.Ctx, update *bot.Update) (bool, *string, error) {
	price, err := controller.getTicketPrice(c)
	if err != nil {
		HandlerPrintf(c, "Failed to get ticket price - %v", err)
		return false, nil, err
	}

	if price == nil {
		message := "Ticket price is not set"
		return false, &message, nil
	}

	if update.PreCheckoutQuery.Currency != price.Currency {
		message := "Incorrect currency"
		return false, &message, nil
	}

	if update.PreCheckoutQuery.TotalAmount != price.Price {
		message := "Incorrect total price"
		return false, &message, nil
	}

	ticketCode, err := uuid.Parse(update.PreCheckoutQuery.InvoicePayload)
	if err != nil {
		message := "Incorrect ticket code"
		return false, &message, nil
	}

	ticket, err := controller.TicketRepository.GetTicket(ticketCode.String())
	if err != nil {
		HandlerPrintf(c, "Failed to get ticket  - %v", err)
		return false, nil, err
	}

	if ticket != nil {
		message := "Ticket already purchased"
		return false, &message, nil
	}

	return true, nil, nil
}

type TicketPrice struct {
	Currency string
	Price    int64
}

func (controller *BotController) getTicketPrice(c *fiber.Ctx) (*TicketPrice, error) {
	result := TicketPrice{}

	currency, err := controller.ConfigRepository.GetValue(TICKET_CURRENCY_KEY)
	if err != nil {
		return nil, err
	}

	if currency == nil {
		HandlerPrintf(c, "Ticket currency key not found")
		return nil, nil
	}

	priceString, err := controller.ConfigRepository.GetValue(TICKET_PRICE_KEY)
	if err != nil {
		return nil, err
	}

	if priceString == nil {
		HandlerPrintf(c, "Ticket price key not found")
		return nil, nil
	}

	price, err := strconv.ParseInt(*priceString, 10, 64)
	if err != nil {
		return nil, err
	}

	result.Currency = *currency
	result.Price = price

	return &result, nil
}

func (controller *BotController) buildWelcomeMessage() (string, bot.SendMessageOptions) {
	message := "Let's start the tour!🎧\nPlease choose an option below to proceed"
	appURL := controller.WebAppURL
	callbackQuery := BUY_TICKET_QUERY
	opts := bot.SendMessageOptions{
		InlineKeyboard: &bot.InlineKeyboardMarkup{
			Markup: [][]bot.InlineKeyboardButton{{
				{Text: "Start the tour", WebAppURL: &appURL},
			}, {
				{Text: "Buy a ticket", CallbackData: &callbackQuery},
			}},
		},
	}

	return message, opts
}

func (controller *BotController) buildPurchaseMessage(ticketCode uuid.UUID) (string, bot.SendMessageOptions) {
	ticketString := ticketCode.String()
	message := "Thank you for your purchase!\nYour ticket number: " + ticketString + "\nPlease tap the button below to proceed"
	appURL := controller.WebAppURL + "?ticket=" + ticketString
	opts := bot.SendMessageOptions{
		InlineKeyboard: &bot.InlineKeyboardMarkup{
			Markup: [][]bot.InlineKeyboardButton{{
				{Text: "Start the tour", WebAppURL: &appURL},
			}},
		},
	}

	return message, opts
}

func (controller *BotController) buildPaymentsDisabledMessage() (string, bot.SendMessageOptions) {
	message := "Sorry, payments are currently not available. Please try again later."

	return message, bot.SendMessageOptions{}
}
