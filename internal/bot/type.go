package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// INFO: CustomType for message for log time respose create
type MessageWithTime struct {
	Msg        *tgbotapi.MessageConfig
	EditMesage *tgbotapi.EditMessageTextConfig
	WorkTime   time.Time
}
