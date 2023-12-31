package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/provider/bot"
	"github.com/st-matskevich/audio-guide-bot/api/provider/translation"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
)

const BUY_TICKET_QUERY = "buy_ticket"
const TICKET_CURRENCY_KEY = "TICKET_CURRENCY"
const TICKET_PRICE_KEY = "TICKET_PRICE"

type BotController struct {
	WebAppURL           string
	BotProvider         bot.BotProvider
	TranslationProvider translation.TranslationProvider
	TicketRepository    repository.TicketRepository
	ConfigRepository    repository.ConfigRepository
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
		HandlerPrintf(c, LOG_ERROR, "Failed to parse input", "error", err)
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Failed to parse input")
	}

	if update.CallbackQuery != nil {
		HandlerPrintf(c, LOG_INFO, "Received callback query")
		return controller.HandleBotCallback(c, &update)
	}

	if update.PreCheckoutQuery != nil {
		HandlerPrintf(c, LOG_INFO, "Received pre-checkout query")
		return controller.HandleBotPreCheckout(c, &update)
	}

	if update.Message != nil {
		HandlerPrintf(c, LOG_INFO, "Received message")
		return controller.HandleBotMessage(c, &update)
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotMessage(c *fiber.Ctx, update *bot.Update) error {
	locale := update.Message.From.LanguageCode

	if update.Message.SuccessfulPayment != nil {
		HandlerPrintf(c, LOG_INFO, "Message type is successful payment")
		ticketCode, err := uuid.Parse(update.Message.SuccessfulPayment.InvoicePayload)
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to parse bot payment payload", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to parse bot payment payload")
		}

		if err = controller.TicketRepository.CreateTicket(ticketCode.String()); err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to register ticket in DB", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to register ticket in DB")
		}

		message, options, err := controller.buildPurchaseMessage(locale, ticketCode.String())
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to prepare message", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
		}

		HandlerPrintf(c, LOG_INFO, "Created ticket, responding with payment confirmation", "ticket", ticketCode.String())
		if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to send bot message", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	HandlerPrintf(c, LOG_INFO, "Responding with welcome message")
	message, options, err := controller.buildWelcomeMessage(locale)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to prepare message", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
	}

	if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to send bot message", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotCallback(c *fiber.Ctx, update *bot.Update) error {
	locale := update.CallbackQuery.From.LanguageCode

	if update.CallbackQuery.Data == "" {
		HandlerPrintf(c, LOG_ERROR, "Bot update didn't include a callback data")
		return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback data")
	}

	if err := controller.BotProvider.AnswerCallbackQuery(update.CallbackQuery.Id); err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to answer callback query", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer callback query")
	}

	if update.CallbackQuery.Data == BUY_TICKET_QUERY {
		HandlerPrintf(c, LOG_INFO, "Callback query is BUY_TICKET_QUERY")
		if update.CallbackQuery.Message == nil {
			HandlerPrintf(c, LOG_ERROR, "Bot update didn't include a callback message")
			return HandlerSendFailure(c, fiber.StatusBadRequest, "Bot update didn't include a callback message")
		}

		price, err := controller.getTicketPrice(c)
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to get ticket price", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to get ticket price")
		} else if price == nil {
			HandlerPrintf(c, LOG_INFO, "Ticket price is not set, responding with disabled payments message")
			message, options, err := controller.buildPaymentsDisabledMessage(locale)
			if err != nil {
				HandlerPrintf(c, LOG_ERROR, "Failed to prepare message", "error", err)
				return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
			}

			if err := controller.BotProvider.SendMessage(update.CallbackQuery.Message.Chat.Id, message, options); err != nil {
				HandlerPrintf(c, LOG_ERROR, "Failed to send bot message", "error", err)
				return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
			}

			return HandlerSendSuccess(c, fiber.StatusOK, nil)
		}

		invoice, err := controller.buildInvoiceData(locale, *price)
		if err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to prepare invoice", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare invoice")
		}

		ticketCode := uuid.New()
		HandlerPrintf(c, LOG_INFO, "Responding with invoice for ticket", "ticket", ticketCode.String())
		if err := controller.BotProvider.SendInvoice(update.CallbackQuery.Message.Chat.Id, invoice.Title, invoice.Description, ticketCode.String(), invoice.Price, invoice.Options); err != nil {
			HandlerPrintf(c, LOG_ERROR, "Failed to send bot invoice", "error", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot invoice")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotPreCheckout(c *fiber.Ctx, update *bot.Update) error {
	locale := update.PreCheckoutQuery.From.LanguageCode
	acceptCheckout, errorMessage, err := controller.validatePreCheckoutQuery(c, update, locale)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to validate pre-checkout query", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to validate pre-checkout query")
	}

	HandlerPrintf(c, LOG_INFO, "Accepting pre-checkout", "accept", acceptCheckout)
	if err := controller.BotProvider.AnswerPreCheckoutQuery(update.PreCheckoutQuery.Id, acceptCheckout, bot.AnswerPreCheckoutQueryOptions{ErrorMessage: errorMessage}); err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to answer pre-checkout query", "error", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer pre-checkout query")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) validatePreCheckoutQuery(c *fiber.Ctx, update *bot.Update, locale string) (bool, *string, error) {
	price, err := controller.getTicketPrice(c)
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to get ticket price", "error", err)
		return false, nil, err
	}

	if price == nil {
		HandlerPrintf(c, LOG_ERROR, "Ticket price is not set")
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_PRICE_NOT_SET", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	if update.PreCheckoutQuery.Currency != price.Currency {
		HandlerPrintf(c, LOG_WARNING, "Pre-checkout currency is not correct")
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_CURRENCY", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	if update.PreCheckoutQuery.TotalAmount != price.Price {
		HandlerPrintf(c, LOG_WARNING, "Pre-checkout price is not correct")
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_PRICE", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	ticketCode, err := uuid.Parse(update.PreCheckoutQuery.InvoicePayload)
	if err != nil {
		HandlerPrintf(c, LOG_WARNING, "Pre-checkout payload is not correct")
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_TICKET", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	ticket, err := controller.TicketRepository.GetTicket(ticketCode.String())
	if err != nil {
		HandlerPrintf(c, LOG_ERROR, "Failed to get invoice ticket", "error", err)
		return false, nil, err
	}

	if ticket != nil {
		HandlerPrintf(c, LOG_WARNING, "Invoice ticket is already sold", "ticket", ticketCode.String())
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_TICKET_SOLD", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	return true, nil, nil
}

type TicketPrice struct {
	Currency string
	Price    int64
}

func (controller *BotController) getTicketPrice(c *fiber.Ctx) (*TicketPrice, error) {
	currency, err := controller.ConfigRepository.GetValue(TICKET_CURRENCY_KEY)
	if err != nil {
		return nil, err
	}

	if currency == nil {
		HandlerPrintf(c, LOG_ERROR, "Ticket currency key not found")
		return nil, nil
	}

	priceString, err := controller.ConfigRepository.GetValue(TICKET_PRICE_KEY)
	if err != nil {
		return nil, err
	}

	if priceString == nil {
		HandlerPrintf(c, LOG_ERROR, "Ticket price key not found")
		return nil, nil
	}

	price, err := strconv.ParseInt(*priceString, 10, 64)
	if err != nil {
		return nil, err
	}

	result := TicketPrice{
		Currency: *currency,
		Price:    price,
	}

	return &result, nil
}

func (controller *BotController) buildWelcomeMessage(locale string) (string, bot.SendMessageOptions, error) {
	message, err := controller.TranslationProvider.TranslateMessage("MESSAGE_WELCOME", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	startText, err := controller.TranslationProvider.TranslateMessage("BUTTON_START_TOUR", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	buyTicketText, err := controller.TranslationProvider.TranslateMessage("BUTTON_BUY_TICKET", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	appURL := controller.WebAppURL
	callbackQuery := BUY_TICKET_QUERY
	opts := bot.SendMessageOptions{
		InlineKeyboard: &bot.InlineKeyboardMarkup{
			Markup: [][]bot.InlineKeyboardButton{{
				{Text: startText, WebAppURL: &appURL},
			}, {
				{Text: buyTicketText, CallbackData: &callbackQuery},
			}},
		},
	}

	return message, opts, nil
}

func (controller *BotController) buildPurchaseMessage(locale string, ticketCode string) (string, bot.SendMessageOptions, error) {
	message, err := controller.TranslationProvider.TranslateMessage("MESSAGE_PURCHASED_TICKET", locale, translation.TemplateData{"TICKET_CODE": ticketCode})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	startText, err := controller.TranslationProvider.TranslateMessage("BUTTON_START_TOUR", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	appURL := controller.WebAppURL + "?ticket=" + ticketCode
	opts := bot.SendMessageOptions{
		InlineKeyboard: &bot.InlineKeyboardMarkup{
			Markup: [][]bot.InlineKeyboardButton{{
				{Text: startText, WebAppURL: &appURL},
			}},
		},
	}

	return message, opts, nil
}

func (controller *BotController) buildPaymentsDisabledMessage(locale string) (string, bot.SendMessageOptions, error) {
	message, err := controller.TranslationProvider.TranslateMessage("MESSAGE_PAYMENTS_NOT_AVAILABLE", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	return message, bot.SendMessageOptions{}, nil
}

type InvoiceData struct {
	Title       string
	Description string
	Price       bot.InvoicePrice
	Options     bot.SendInvoiceOptions
}

func (controller *BotController) buildInvoiceData(locale string, price TicketPrice) (InvoiceData, error) {
	title, err := controller.TranslationProvider.TranslateMessage("PAYMENT_TICKET_TITLE", locale, translation.TemplateData{})
	if err != nil {
		return InvoiceData{}, err
	}

	description, err := controller.TranslationProvider.TranslateMessage("PAYMENT_TICKET_DESCRIPTION", locale, translation.TemplateData{})
	if err != nil {
		return InvoiceData{}, err
	}

	priceLabel, err := controller.TranslationProvider.TranslateMessage("PAYMENT_TICKET_PRICE_PART_PRICE", locale, translation.TemplateData{})
	if err != nil {
		return InvoiceData{}, err
	}

	payText, err := controller.TranslationProvider.TranslateMessage("BUTTON_PAY", locale, translation.TemplateData{})
	if err != nil {
		return InvoiceData{}, err
	}

	pay := true
	opts := bot.SendInvoiceOptions{
		InlineKeyboard: &bot.InlineKeyboardMarkup{
			Markup: [][]bot.InlineKeyboardButton{{
				{Text: payText, Pay: &pay},
			}},
		},
	}

	data := InvoiceData{
		Title:       title,
		Description: description,
		Price: bot.InvoicePrice{
			Currency: price.Currency,
			Parts:    []bot.PricePart{{Label: priceLabel, Amount: price.Price}},
		},
		Options: opts,
	}

	return data, nil
}
