package telegram

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler for /start command
func (h *Hub) cmdStrart(chatID int64, userName string) (msg *tgbotapi.MessageConfig, err error) {
	// obj user
	u := &store.User{
		ChatID:   chatID,
		UserName: &userName,
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
	// User have not NickName
	if userName == "" {
		msg = h.createMessage(chatID, NotNickNameUserMsg)
	} else {
		// create msg
		msg = h.createMessage(chatID, cmdStart)
	}
	return msg, nil
}

func (h *Hub) cmdGetMenu(chatID int64) *tgbotapi.MessageConfig {
	msg := h.createMessage(chatID, cmdMenu)
	msg.ReplyMarkup = createInlineGetChoiceList()
	return msg
}

func (h *Hub) cmdHelp(chatID int64) *tgbotapi.MessageConfig {
	msg := h.createMessage(chatID, cmdHelpMessage)
	return msg
}

// Handler for default message
func (h *Hub) cmdDefault(chatID int64) *tgbotapi.MessageConfig {
	msg := h.createMessage(chatID, "I don`t know this command")
	return msg
}
