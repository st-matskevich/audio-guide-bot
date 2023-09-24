package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type BotInteractor interface {
	SendMessage(chatId int64, text string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error)
}
