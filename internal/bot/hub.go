package telegram

import (
	"context"
	"errors"
	"strings"

	manager "github.com/Negat1v9/telegram-bot-orders/internal"
	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	NoCallbackDataError = errors.New("Not exists handler for CallBack")
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
	// command message
	answer, err = h.isCommand(text, msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}
	// Message is Forward
	answer, err = h.isForwardMessage(msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}
	// only text message
	answer, err = h.isMessage(text, msg)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		return answer, nil
	}
	return h.cmdDefault(msg.Chat.ID), nil
}

func (h *Hub) CallBackUpdate(cbq tgbotapi.CallbackQuery) (*tgbotapi.MessageConfig, error) {
	switch {
	case isGetProductList(cbq.Data):
		listID, listName := parseIDNameList(cbq.Data)
		msg, err := h.getProductList(cbq.From.ID, listID, listName)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case isAddNewProduct(cbq.Data):
		listName := parseNameListFromProductAction(cbq.Data)
		msg := h.createMessage(cbq.From.ID, addNewProductMessage+listName)
		return msg, nil
		// case isCreateAddProduct(cbq.Data):
		// 	listName := parseNameListFromProductAction(cbq.Data)
		// 	msg := h.
	}
	return nil, NoCallbackDataError
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
		msg, err := h.getListName(msgInfo.From.ID, msgInfo.Chat.ID)
		if err != nil {
			return nil, err
		}
		return msg, nil
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
	case isAddNewProductForward(text):
		listName := parseNameList(text)
		listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
		if err != nil {
			return nil, err
		}
		products := parseStringToProducts(msg.Text, listID)

		msg, err := h.addNewProduct(msg.Chat.ID, products, listName)
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
