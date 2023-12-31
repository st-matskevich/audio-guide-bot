package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type Update gotgbot.Update

type InlineKeyboardButton struct {
	Text         string
	URL          *string
	WebAppURL    *string
	CallbackData *string
	Pay          *bool
}

type InlineKeyboardMarkup struct {
	Markup [][]InlineKeyboardButton
}

type SendMessageOptions struct {
	InlineKeyboard *InlineKeyboardMarkup
}

type SendInvoiceOptions struct {
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

type AnswerPreCheckoutQueryOptions struct {
	ErrorMessage *string
}

type BotProvider interface {
	SendMessage(chatID int64, text string, options SendMessageOptions) error
	AnswerCallbackQuery(queryID string) error
	AnswerPreCheckoutQuery(queryID string, ok bool, options AnswerPreCheckoutQueryOptions) error
	SendInvoice(chatID int64, title string, description string, payload string, price InvoicePrice, options SendInvoiceOptions) error
}
