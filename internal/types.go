package manager

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Manager interface {
	MessageUpdate(msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}
