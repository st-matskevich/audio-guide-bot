package bot

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type TelegramBotInteractor struct {
	PaymentsToken string
	Bot           *gotgbot.Bot
}

func (interactor *TelegramBotInteractor) SendMessage(chatID int64, text string, options SendMessageOptions) error {
	opts := &gotgbot.SendMessageOpts{}

	if options.InlineKeyboard != nil {
		markup := [][]gotgbot.InlineKeyboardButton{}
		for _, row := range options.InlineKeyboard.Markup {
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

				inlineRow = append(inlineRow, inlineButton)
			}
			markup = append(markup, inlineRow)
		}

		opts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: markup,
		}
	}

	if _, err := interactor.Bot.SendMessage(chatID, text, opts); err != nil {
		return err
	}

	return nil
}

func (interactor *TelegramBotInteractor) AnswerCallbackQuery(queryID string) error {
	result, err := interactor.Bot.AnswerCallbackQuery(queryID, nil)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("AnswerCallbackQuery returned false")
	}

	return nil
}

func (interactor *TelegramBotInteractor) AnswerPreCheckoutQuery(queryID string, ok bool) error {
	result, err := interactor.Bot.AnswerPreCheckoutQuery(queryID, ok, nil)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("AnswerPreCheckoutQuery returned false")
	}

	return nil
}

func (interactor *TelegramBotInteractor) SendInvoice(chatID int64, title string, description string, payload string, price InvoicePrice) error {
	labeledPrice := []gotgbot.LabeledPrice{}
	for _, part := range price.Parts {
		labeledPrice = append(labeledPrice, gotgbot.LabeledPrice{Label: part.Label, Amount: part.Amount})
	}

	if _, err := interactor.Bot.SendInvoice(chatID, title, description, payload, interactor.PaymentsToken, price.Currency, labeledPrice, nil); err != nil {
		return err
	}

	return nil
}

func CreateTelegramBotInteractor(botToken string, paymentsToken string) (BotInteractor, error) {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		return nil, err
	}

	interactor := TelegramBotInteractor{
		Bot:           bot,
		PaymentsToken: paymentsToken,
	}

	return &interactor, nil
}
