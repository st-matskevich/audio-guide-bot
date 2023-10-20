package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/st-matskevich/audio-guide-bot/api/bot"
	"github.com/st-matskevich/audio-guide-bot/api/repository"
	"github.com/st-matskevich/audio-guide-bot/api/translation"
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
	locale := update.Message.From.LanguageCode

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

		message, options, err := controller.buildPurchaseMessage(locale, ticketCode)
		if err != nil {
			HandlerPrintf(c, "Failed to prepare message - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
		}

		if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
			HandlerPrintf(c, "Failed to send bot message - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
		}

		return HandlerSendSuccess(c, fiber.StatusOK, nil)
	}

	message, options, err := controller.buildWelcomeMessage(locale)
	if err != nil {
		HandlerPrintf(c, "Failed to prepare message - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
	}

	if err := controller.BotProvider.SendMessage(update.Message.Chat.Id, message, options); err != nil {
		HandlerPrintf(c, "Failed to send bot message - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) HandleBotCallback(c *fiber.Ctx, update *bot.Update) error {
	locale := update.CallbackQuery.From.LanguageCode

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

			message, options, err := controller.buildPaymentsDisabledMessage(locale)
			if err != nil {
				HandlerPrintf(c, "Failed to prepare message - %v", err)
				return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare message")
			}

			if err := controller.BotProvider.SendMessage(update.CallbackQuery.Message.Chat.Id, message, options); err != nil {
				HandlerPrintf(c, "Failed to send bot message - %v", err)
				return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to send bot message")
			}

			return HandlerSendSuccess(c, fiber.StatusOK, nil)
		}

		invoice, err := controller.buildInvoiceData(locale, *price)
		if err != nil {
			HandlerPrintf(c, "Failed to prepare invoice - %v", err)
			return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to prepare invoice")
		}

		ticketCode := uuid.New()
		if err := controller.BotProvider.SendInvoice(update.CallbackQuery.Message.Chat.Id, invoice.Title, invoice.Description, ticketCode.String(), invoice.Price, invoice.Options); err != nil {
			HandlerPrintf(c, "Failed to send bot invoice - %v", err)
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
		HandlerPrintf(c, "Failed to validate pre-checkout query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to validate pre-checkout query")
	}

	if err := controller.BotProvider.AnswerPreCheckoutQuery(update.PreCheckoutQuery.Id, acceptCheckout, bot.AnswerPreCheckoutQueryOptions{ErrorMessage: errorMessage}); err != nil {
		HandlerPrintf(c, "Failed to answer pre-checkout query - %v", err)
		return HandlerSendError(c, fiber.StatusInternalServerError, "Failed to answer pre-checkout query")
	}

	return HandlerSendSuccess(c, fiber.StatusOK, nil)
}

func (controller *BotController) validatePreCheckoutQuery(c *fiber.Ctx, update *bot.Update, locale string) (bool, *string, error) {
	price, err := controller.getTicketPrice(c)
	if err != nil {
		HandlerPrintf(c, "Failed to get ticket price - %v", err)
		return false, nil, err
	}

	if price == nil {
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_PRICE_NOT_SET", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	if update.PreCheckoutQuery.Currency != price.Currency {
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_CURRENCY", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	if update.PreCheckoutQuery.TotalAmount != price.Price {
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_PRICE", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	ticketCode, err := uuid.Parse(update.PreCheckoutQuery.InvoicePayload)
	if err != nil {
		message, err := controller.TranslationProvider.TranslateMessage("PAYMENT_FAIL_INVALID_TICKET", locale, translation.TemplateData{})
		if err != nil {
			return false, nil, err
		}
		return false, &message, nil
	}

	ticket, err := controller.TicketRepository.GetTicket(ticketCode.String())
	if err != nil {
		HandlerPrintf(c, "Failed to get ticket  - %v", err)
		return false, nil, err
	}

	if ticket != nil {
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

func (controller *BotController) buildPurchaseMessage(locale string, ticketCode uuid.UUID) (string, bot.SendMessageOptions, error) {
	ticketString := ticketCode.String()
	message, err := controller.TranslationProvider.TranslateMessage("MESSAGE_PURCHASED_TICKET", locale, translation.TemplateData{"TICKET_CODE": ticketString})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	startText, err := controller.TranslationProvider.TranslateMessage("BUTTON_START_TOUR", locale, translation.TemplateData{})
	if err != nil {
		return "", bot.SendMessageOptions{}, err
	}

	appURL := controller.WebAppURL + "?ticket=" + ticketString
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
