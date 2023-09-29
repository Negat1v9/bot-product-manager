package telegram

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler for /start command
func (h *Hub) cmdStrart(userName string, chatID int64) (*tgbotapi.MessageConfig, error) {
	// obj user
	u := &store.User{
		ChatID:   int(chatID),
		UserName: userName,
	}
	// To confirm user in db or not
	isExist, err := h.db.User().IsExist(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	// if current user is new -> create one
	if !(isExist) {
		err = h.db.User().Add(context.TODO(), u)
	}
	if err != nil {
		return nil, err
	}
	// create msg
	msg := h.createMessage(chatID, "Hello from start command!")
	msg.ReplyMarkup = menuKeyboard
	return msg, nil
}

func (h *Hub) cmdHelp(chatID int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, cmdHelpMessage)
	return &msg
}

// Handler for default message
func (h *Hub) cmdDefault(chatID int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Default command")
	return &msg
}
