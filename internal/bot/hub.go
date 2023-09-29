package telegram

import (
	"strings"

	manager "github.com/Negat1v9/telegram-bot-orders/internal"
	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Hub struct {
	db store.Store
}

func NewHub(db store.Store) manager.Manager {
	return &Hub{
		db: db,
	}
}

func (h *Hub) MessageUpdate(msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	text := toLowerCase(msg.Text)
	var answer *tgbotapi.MessageConfig
	var err error

	answer, err = h.isCommand(text, msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}

	answer, err = h.isForwardMessage(msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}

	answer, err = h.isMessage(text, msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}
	return h.cmdDefault(msg.Chat.ID), nil
}

func (h *Hub) isCommand(text string, msgInfo *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	switch text {
	case "/start":
		msg, err := h.cmdStrart(msgInfo.From.UserName, msgInfo.Chat.ID)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case "/help":
		msg := h.cmdHelp(msgInfo.Chat.ID)
		return msg, nil
	}
	return nil, nil
}

func (h *Hub) isMessage(text string, msgInfo *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	switch {
	case isCreateList(text):
		msg := h.answerToCreateList(msgInfo.From.ID)
		return msg, nil
	case isSelectUserList(text):
		msg, err := h.selectList(msgInfo.From.ID, msgInfo.Chat.ID)
		if err != nil {
			return nil, err
		}
		return msg, nil
		// case isAddCommand(text):
		// 	// TODO: add func to add product in list
	}
	return nil, nil
}

func (h *Hub) isForwardMessage(msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	if msg.ReplyToMessage == nil {
		return nil, nil
	}
	text := msg.ReplyToMessage.Text
	switch {
	case isCreateNameForward(text):
		list := &store.ProductList{
			OwnerID: msg.From.ID,
			Name:    msg.Text,
		}
		msg, err := h.createList(msg.Chat.ID, list)
		if err != nil {
			return nil, err
		}
		return msg, nil
	default:
		return nil, nil

	}
}

func (h *Hub) createMessage(ChatId int64, text string) *tgbotapi.MessageConfig {
	msgCongig := tgbotapi.NewMessage(ChatId, text)
	return &msgCongig
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}
