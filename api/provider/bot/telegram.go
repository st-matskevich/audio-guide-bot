package bot

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type TelegramBotProvider struct {
	PaymentsToken string
	Bot           *gotgbot.Bot
}

func (interactor *TelegramBotProvider) SendMessage(chatID int64, text string, options SendMessageOptions) error {
	opts := &gotgbot.SendMessageOpts{}

	if options.InlineKeyboard != nil {
		opts.ReplyMarkup = interactor.buildKeyboard(*options.InlineKeyboard)
	}

	if _, err := interactor.Bot.SendMessage(chatID, text, opts); err != nil {
		return err
	}

	return nil
}

func (interactor *TelegramBotProvider) AnswerCallbackQuery(queryID string) error {
	result, err := interactor.Bot.AnswerCallbackQuery(queryID, nil)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("AnswerCallbackQuery returned false")
	}

	return nil
}

func (interactor *TelegramBotProvider) AnswerPreCheckoutQuery(queryID string, ok bool, options AnswerPreCheckoutQueryOptions) error {
	opts := &gotgbot.AnswerPreCheckoutQueryOpts{}
	if options.ErrorMessage != nil {
		opts.ErrorMessage = *options.ErrorMessage
	}

	result, err := interactor.Bot.AnswerPreCheckoutQuery(queryID, ok, opts)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("AnswerPreCheckoutQuery returned false")
	}

	return nil
}

func (interactor *TelegramBotProvider) SendInvoice(chatID int64, title string, description string, payload string, price InvoicePrice, options SendInvoiceOptions) error {
	labeledPrice := []gotgbot.LabeledPrice{}
	for _, part := range price.Parts {
		labeledPrice = append(labeledPrice, gotgbot.LabeledPrice{Label: part.Label, Amount: part.Amount})
	}

	opts := &gotgbot.SendInvoiceOpts{}
	if options.InlineKeyboard != nil {
		opts.ReplyMarkup = interactor.buildKeyboard(*options.InlineKeyboard)
	}

	if _, err := interactor.Bot.SendInvoice(chatID, title, description, payload, interactor.PaymentsToken, price.Currency, labeledPrice, opts); err != nil {
		return err
	}

	return nil
}

func (interactor *TelegramBotProvider) buildKeyboard(keyboard InlineKeyboardMarkup) gotgbot.InlineKeyboardMarkup {
	markup := [][]gotgbot.InlineKeyboardButton{}
	for _, row := range keyboard.Markup {
		inlineRow := []gotgbot.InlineKeyboardButton{}
		for _, button := range row {
			inlineButton := gotgbot.InlineKeyboardButton{
				Text: button.Text,
			}

			if button.URL != nil {
				inlineButton.Url = *button.URL
			}

			if button.WebAppURL != nil {
				inlineButton.WebApp = &gotgbot.WebAppInfo{Url: *button.WebAppURL}
			}

			if button.CallbackData != nil {
				inlineButton.CallbackData = *button.CallbackData
			}

			if button.Pay != nil {
				inlineButton.Pay = *button.Pay
			}

			inlineRow = append(inlineRow, inlineButton)
		}
		markup = append(markup, inlineRow)
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}

func CreateTelegramBotProvider(botToken string, paymentsToken string) (BotProvider, error) {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		return nil, err
	}

	interactor := TelegramBotProvider{
		Bot:           bot,
		PaymentsToken: paymentsToken,
	}

	return &interactor, nil
}
