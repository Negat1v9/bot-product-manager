package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	manager "github.com/Negat1v9/telegram-bot-orders/internal"
	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CommandUpdate = iota
	TextUpdate
	ForwardMessageUpdate
	CallBackUpdate
	InvalidTypeUpdate
)

var (
	NoCallbackDataError = errors.New("Not exists handler for CallBack")
)

type Hub struct {
	db       store.Store
	response chan<- *MessageWithTime
}

func NewHub(db store.Store, resCh chan<- *MessageWithTime) manager.Manager {
	return &Hub{
		db:       db,
		response: resCh,
	}
}

func (h *Hub) MessageUpdate(msg *tg.Message, timeStart time.Time) (err error) {
	text := toLowerCase(msg.Text)

	var res *tg.MessageConfig

	typeUpdate := h.getTypeMessage(text, msg)

	switch typeUpdate {
	case CommandUpdate:
		res, err = h.isCommand(text, msg)

	case ForwardMessageUpdate:
		res, err = h.isForwardMessage(msg)

	case TextUpdate:
		res, err = h.isMessage(text, msg)

	}
	if err != nil {
		return err
	}
	h.response <- &MessageWithTime{res, timeStart}
	return nil
}

func (h *Hub) getTypeMessage(text string, msg *tg.Message) int {
	if isCommandUpdate(text) {
		return CommandUpdate
	}
	if msg.ReplyToMessage != nil {
		return ForwardMessageUpdate
	}
	return TextUpdate
}

func (h *Hub) CallBackUpdate(cbq *tg.CallbackQuery, timeStart time.Time) error {
	var res *tg.MessageConfig
	var err error
	switch {

	case isGetProductList(cbq.Data):
		listID, listName := parseIDName(cbq.Data)
		res, err = h.getProductList(cbq.From.ID, listID, listName)

	case isAddNewProduct(cbq.Data):
		listName := parseNameListFromProductAction(cbq.Data)
		res = h.createMessage(cbq.From.ID, addNewProductMessage+listName)

	case isGetGroupLists(cbq.Data):
		groupID := parseGroupID(cbq.Data)
		res, err = h.GetGroupLists(cbq.From.ID, groupID)

	case isCreateGroupList(cbq.Data):
		groupID := parseGroupID(cbq.Data)
		res, err = h.createMessageCreateGroupList(cbq.From.ID, groupID)

	case isCompliteProductList(cbq.Data):
		listName := parseNameListFromProductAction(cbq.Data)
		res, err = h.compliteProductList(cbq.From.ID, listName)

	case isGetUsersForDelGroup(cbq.Data):
		groupID := parseGroupID(cbq.Data)
		res, err = h.getUserForDeleteFrGr(cbq.From.ID, groupID)
	case isAddNewUserGroup(cbq.Data):
		groupID := parseGroupID(cbq.Data)
		res, err = h.createMessageForInviteUser(cbq.From.ID, groupID)

	case isUserReadyInvite(cbq.Data):
		userID, groupID := parseCallBackGroupActions(cbq.Data)
		fmt.Println("hub:", cbq.Data)
		res, err = h.userReadyJoinGroup(cbq.From.ID, userID, groupID)

	case isDeleteUserFromGroup(cbq.Data):
		userID, groupID := parseCallBackGroupActions(cbq.Data)
		res, err = h.deleteUserFromGroup(cbq.From.ID, userID, groupID)
	}

	if err != nil {
		return err
	}
	if res != nil {
		h.response <- &MessageWithTime{res, timeStart}
		return nil
	}
	return NoCallbackDataError
}

func (h *Hub) isCommand(text string, msgInfo *tg.Message) (*tg.MessageConfig, error) {
	switch text {
	case "/start":
		msg, err := h.cmdStrart(msgInfo.From.UserName, msgInfo.From.ID)
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

func (h *Hub) isMessage(text string, msgInfo *tg.Message) (*tg.MessageConfig, error) {
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
	case isGetUserGroup(text):
		msg, err := h.GetAllUserGroup(msgInfo.Chat.ID, msgInfo.From.ID)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case isCreateGroup(text):
		msg := h.createMessageForNewGroup(msgInfo.Chat.ID)
		return msg, nil
	}
	return nil, nil
}

func (h *Hub) isForwardMessage(msg *tg.Message) (*tg.MessageConfig, error) {
	text := msg.ReplyToMessage.Text
	var res *tg.MessageConfig
	var err error
	switch {
	case isCreateNameForward(text):
		list := &store.ProductList{
			OwnerID: &msg.From.ID,
			Name:    &msg.Text,
		}
		res, err = h.createList(msg.Chat.ID, list)

	case isAddNewProductForward(text):
		listName := parseNameListForAddProd(text)
		// FIXME: function for getting listID put in addNewProduct()
		listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
		if err != nil {
			return nil, err
		}
		products := parseStringToProducts(msg.Text, listID)

		res, err = h.addNewProduct(msg.Chat.ID, products, listName)

	case isCreateNewGroupForward(text):
		managerGroup := &store.GroupInfo{
			OwnerID:   msg.From.ID,
			GroupName: msg.Text,
		}
		res, err = h.createNewGroup(msg.Chat.ID, managerGroup)

	case isCreateGroupListForward(text):
		newListName := parseGroupListName(msg.Text)
		groupName := parseGroupListName(text)
		res, err = h.createGroupList(msg.From.ID, newListName, groupName)

	case isSendInviteToNewUser(text):
		groupName := parseNameGroupAddUser(text)
		newUserName := parseUserNickNameForAddGroup(msg.Text)
		res, err = h.inviteNewUser(msg.From.ID, newUserName, groupName)
	}
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return nil, nil
}

func (h *Hub) createMessage(ChatId int64, text string) *tg.MessageConfig {
	msgCongig := tg.NewMessage(ChatId, text)
	return &msgCongig
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}
