package telegram

import (
	"errors"
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
	response chan<- MessageWithTime
}

func NewHub(db store.Store, resCh chan<- MessageWithTime) manager.Manager {
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
		msg := h.createErrorMessaeg(msg.From.ID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return err
	}
	if res != nil {
		h.response <- MessageWithTime{Msg: res, WorkTime: timeStart}
		return nil
	}
	// send default message if user print something what we don`t know
	h.response <- MessageWithTime{
		Msg:      h.cmdDefault(msg.From.ID),
		WorkTime: timeStart,
	}
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
	var msg *tg.MessageConfig
	var editMsg *tg.EditMessageTextConfig
	var err error
	switch {
	// TODO:
	case isUserChoiceLists(cbq.Data):
		editMsg, err = h.getListName(cbq.From.ID, cbq.Message.MessageID)
	// editMsg, err = h.getListName()
	case isUserChoiceGroupLists(cbq.Data):
		editMsg, err = h.GetAllUserGroup(cbq.From.ID, cbq.Message.MessageID)
	case isGetProductList(cbq.Data):
		data := parseCallBackFewParam(prefixCallBackListProduct, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.getProductList(cbq.From.ID, cbq.Message.MessageID, listID, listName)

	case isAddNewProduct(cbq.Data):
		listName := parseCallBackOneParam(prefixAddProductList, cbq.Data)
		msg = h.createMessage(cbq.From.ID, addNewProductMessage+listName)

	case isGetGroupLists(cbq.Data):
		data := parseCallBackOneParam(prefixCallBackListGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.GetGroupLists(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isEditProductList(cbq.Data):
		groupName := parseCallBackOneParam(prefixChangeList, cbq.Data)
		msg = h.createMessageForEditList(cbq.From.ID, groupName)

	case isCreateGroupList(cbq.Data):
		data := parseCallBackOneParam(prefixCreateGroupList, cbq.Data)
		groupID := convSToI[int](data, 0)
		msg, err = h.createMessageCreateGroupList(cbq.From.ID, groupID)

	case isCompliteProductList(cbq.Data):
		listName := parseCallBackOneParam(prefixCompliteList, cbq.Data)
		msg, err = h.compliteProductList(cbq.From.ID, listName)

	case isGetUsersForDelGroup(cbq.Data):
		data := parseCallBackOneParam(prefixGetUserToDelete, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.getUserForDeleteFrGr(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isAddNewUserGroup(cbq.Data):
		data := parseCallBackOneParam(prefixAddUserGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		msg, err = h.createMessageForInviteUser(cbq.From.ID, groupID)
		// TODO: Make invite and refuse message as edit type, now its not working
	case isUserReadyInvite(cbq.Data):
		data := parseCallBackFewParam(prefixCallBackInsertUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		msg, err = h.userReadyJoinGroup(userID, groupID)

	case isUserRefuseInvite(cbq.Data):
		data := parseCallBackFewParam(prefixCallBackRefuseUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		msg, err = h.userRefuseJoinGroup(userID, cbq.From.UserName, groupID)

	case isDeleteUserFromGroup(cbq.Data):
		data := parseCallBackFewParam(prefixCallBackDelUserFromGr, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		msg, err = h.deleteUserFromGroup(cbq.From.ID, userID, groupID)
	}
	// TODO: Returning default message if all is nil
	if err != nil {
		msg := h.createErrorMessaeg(cbq.From.ID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return err
	}
	h.response <- MessageWithTime{Msg: msg, EditMesage: editMsg, WorkTime: timeStart}
	return nil
}

func (h *Hub) isCommand(text string, msgInfo *tg.Message) (*tg.MessageConfig, error) {
	switch text {
	case "/start":
		msg, err := h.cmdStrart(msgInfo.From.ID, msgInfo.From.UserName)
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
	case isSelectList(text):
		msg := h.getChoiceLists(msgInfo.From.ID)
		// msg, err := h.getListName(msgInfo.From.ID, msgInfo.Chat.ID)
		// if err != nil {
		// return nil, err
		// }
		return msg, nil
	// case isGetUserGroup(text):
	// 	msg, err := h.GetAllUserGroup(msgInfo.Chat.ID, msgInfo.From.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return msg, nil
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

		res, err = h.addNewProduct(msg.Chat.ID, msg.Text, listName)

	case isEditListForward(text):
		listName := parseListNameEditList(text)
		indexToDelete := parseIndexEditProduct(msg.Text)
		res, err = h.editProductList(msg.From.ID, listName, indexToDelete)

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

func (h *Hub) editMessage(chatID int64, lastmsgDI int, text string) *tg.EditMessageTextConfig {
	msg := tg.NewEditMessageText(chatID, lastmsgDI, text)

	return &msg
}

func (h *Hub) createMessage(ChatId int64, text string) *tg.MessageConfig {
	msgCongig := tg.NewMessage(ChatId, text)
	return &msgCongig
}

func (h *Hub) createErrorMessaeg(chatID int64) *tg.MessageConfig {
	msg := tg.NewMessage(chatID, errorMessage)
	return &msg
}

// NOTE: Make my example
func toLowerCase(s string) string {
	return strings.ToLower(s)
}
