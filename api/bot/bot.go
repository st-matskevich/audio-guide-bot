package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

// TODO: rewrite this provider to hide implementation details from consumers
type BotInteractor interface {
	SendMessage(chatId int64, text string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error)
	AnswerCallbackQuery(callbackQueryId string, opts *gotgbot.AnswerCallbackQueryOpts) (bool, error)
	SendInvoice(chatId int64, title string, description string, payload string, providerToken string, currency string, prices []gotgbot.LabeledPrice, opts *gotgbot.SendInvoiceOpts) (*gotgbot.Message, error)
	AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool, opts *gotgbot.AnswerPreCheckoutQueryOpts) (bool, error)
}
