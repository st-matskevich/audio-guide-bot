package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type Update gotgbot.Update

type InlineKeyboardButton struct {
	Text         string
	URL          *string
	WebAppURL    *string
	CallbackData *string
}

type InlineKeyboardMarkup struct {
	Markup [][]InlineKeyboardButton
}

type SendMessageOptions struct {
	InlineKeyboard *InlineKeyboardMarkup
}

type PricePart struct {
	Label  string
	Amount int64
}

type InvoicePrice struct {
	Currency string
	Parts    []PricePart
}

type BotProvider interface {
	SendMessage(chatID int64, text string, options SendMessageOptions) error
	AnswerCallbackQuery(queryID string) error
	AnswerPreCheckoutQuery(queryID string, ok bool) error
	SendInvoice(chatID int64, title string, description string, payload string, price InvoicePrice) error
}
