package manager

import (
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Manager interface {
	MessageUpdate(msg *tg.Message, t time.Time) error
	CallBackUpdate(cbq *tg.CallbackQuery, t time.Time) error
}
