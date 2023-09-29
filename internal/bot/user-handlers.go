package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: create message with information to user, about he have to forwad its
func (h *Hub) answerToCreateList(ChatID int64) *tgbotapi.MessageConfig {
	msg := h.createMessage(ChatID, answerCreateListMsg)
	return msg
}

// Info: Create list in database, and answer with new list name
func (h *Hub) createList(ChatID int64, list *store.ProductList) (*tgbotapi.MessageConfig, error) {

	err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", list.Name))
	return msg, nil
}

// Info: Select all users lists and create inline keyboard with its
func (h *Hub) selectList(UserID int64, ChatID int64) (*tgbotapi.MessageConfig, error) {
	lists, err := h.db.ProductList().GetAll(context.TODO(), int(UserID))
	if err != nil {
		return nil, err
	}
	keyboard := createInlineMarkup(lists)
	msg := h.createMessage(ChatID, "List of Product-lists")
	msg.ReplyMarkup = keyboard
	return msg, nil
}
