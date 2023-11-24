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
	db        store.Store
	container *Container
	response  chan<- MessageWithTime
}

func NewHub(db store.Store, cont *Container, resCh chan<- MessageWithTime) manager.Manager {
	return &Hub{
		db:        db,
		container: cont,
		response:  resCh,
	}
}

// TODO: For new update create User struct with id and name
func (h *Hub) MessageUpdate(msg *tg.Message, timeStart time.Time) (err error) {
	text := toLowerCase(msg.Text)

	typeUpdate := h.getTypeMessage(text, msg)

	switch typeUpdate {

	case CommandUpdate:
		err = h.isCommand(text, msg, timeStart)

	case ForwardMessageUpdate:
		err = h.isForwardMessage(msg, timeStart)

	}
	if err != nil {
		// TODO: Make it edit type message
		msg := h.createErrorMessaeg(msg.From.ID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return err
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

func (h *Hub) getCallBackType(callBack string) int {
	callBack = clearCallBackData(callBack)
	v, ok := prefixsMap[callBack]
	if !ok {
		return -1
	}
	return v
}

func (h *Hub) CallBackUpdate(cbq *tg.CallbackQuery, timeStart time.Time) error {
	var res = MessageWithTime{WorkTime: timeStart}
	var err error
	CBType := h.getCallBackType(cbq.Data)
	switch CBType {

	case isGetLists:
		res.EditMesage, err = h.getListName(cbq.From.ID, cbq.Message.MessageID)

	case isGetGroupProductList:
		data := parseCallBackFewParam(prefixCallBackGroupProductList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		res.EditMesage, err = h.getGroupList(cbq.From.ID, cbq.Message.MessageID, listID, listName)

	case isGetGroupLists:
		res.EditMesage, err = h.GetAllUserGroup(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateList:
		res.EditMesage = h.createMsgToCreateList(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateGroup:
		res.EditMesage, err = h.createMessageForNewGroup(cbq.From.ID, cbq.Message.MessageID)

	case isGetProductList:
		data := parseCallBackFewParam(prefixCallBackListProduct, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		res.EditMesage, err = h.getProductList(cbq.From.ID, cbq.Message.MessageID, listID, listName)

	case isWantAddNewProduct:
		products := cbq.Message.Text
		listName := parseCallBackOneParam(prefixAddProductList, cbq.Data)
		res.EditMesage = h.wantAddNewProduct(cbq.From.ID, products, listName, cbq.Message.MessageID, false)
		// TODO:
	case isWantAddProductGroupList:
		products := cbq.Message.Text
		listName := parseCallBackOneParam(prefixAddProductGroup, cbq.Data)
		res.EditMesage = h.wantAddNewProduct(cbq.From.ID, products, listName, cbq.Message.MessageID, true)

	case isGetAllGroupLists:
		data := parseCallBackOneParam(prefixCallBackListGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.GetGroupLists(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isWantEditList:
		products := cbq.Message.Text
		groupName := parseCallBackOneParam(prefixChangeList, cbq.Data)
		res.EditMesage = h.createMessageForEditList(cbq.From.ID, products, groupName, cbq.Message.MessageID, false)
		// TODO:
	case isWantEditGroupList:
		products := cbq.Message.Text
		groupName := parseCallBackOneParam(prefixChangeGroupList, cbq.Data)
		res.EditMesage = h.createMessageForEditList(cbq.From.ID, products, groupName, cbq.Message.MessageID, true)

	case isWantCreateGroupList:
		data := parseCallBackOneParam(prefixCreateGroupList, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.createMessageCreateGroupList(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isCompliteSoloList:
		data := parseCallBackFewParam(prefixCompliteSoloList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		res.EditReplyMarkup, err = h.compliteProductList(cbq.From.ID, listName, "", "text", listID, cbq.Message.MessageID)

	case isWantCompliteList:
		products := cbq.Message.Text
		data := parseCallBackFewParam(prefixWantCompliteList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		res.EditMesage, err = h.wantCompliteList(cbq.From.ID, listName, products, listID, cbq.Message.MessageID)

	case isCompliteList:
		data := parseCallBackFewParam(prefixCompliteList, cbq.Data)
		listID, sGroupID, listName := convSToI[int](data[0], 0), data[1], data[2]
		res.EditReplyMarkup, err = h.compliteProductList(cbq.From.ID, listName, sGroupID, cbq.Message.Text, listID, cbq.Message.MessageID)

	case isRestoreProductList:
		data := parseCallBackOneParam(prefixRestoreList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.Msg, err = h.makeListActive(cbq.From.ID, listID)

	case isRestoreGroupList:
		data := parseCallBackFewParam(prefixRestoreGroupList, cbq.Data)
		listID, groupID, listName := convSToI[int](data[0], 0), convSToI[int](data[1], 0), data[2]

		res.Msg, err = h.recoverGroupList(cbq.From.ID, listID, groupID, cbq.Message.MessageID, cbq.Message.Text, listName)
	case isWantMergeList:
		listName := parseCallBackOneParam(prefixToMergeListGroup, cbq.Data)
		res.EditMesage, err = h.getNameGroupMergeList(cbq.From.ID, listName, cbq.Message.MessageID)

	case isGetMainMenu:
		res.EditMesage = h.getMainMenu(cbq.From.ID, cbq.Message.MessageID)

	case isMergeListGroup:
		data := parseCallBackFewParam(prefixMergeListWithGroup, cbq.Data)
		groupID, listID := convSToI[int](data[0], 0), convSToI[int](data[1], 0)
		res.EditMesage, err = h.mergeListWithGroup(cbq.From.ID, groupID, listID, cbq.Message.MessageID)

	case isLeaveGroup:
		data := parseCallBackOneParam(prefixLeaveGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.leaveFromGroup(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isLeaveOwnerGroup:
		data := parseCallBackOneParam(prefixLeaveOwnerGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.leaveAndDeleteGroup(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isGetAllUsersGroup:
		data := parseCallBackOneParam(prefixGetAllUsersGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.getUserFromGroup(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isGetUsersToDelete:
		data := parseCallBackOneParam(prefixGetUserToDelete, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.getUserForDeleteFrGr(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isWantInviteNewUser:
		data := parseCallBackOneParam(prefixAddUserGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.createMessageForInviteUser(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isUserReadyJoinGroup:
		data := parseCallBackFewParam(prefixCallBackInsertUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		res.Msg, err = h.userReadyJoinGroup(userID, groupID)

	case isUserRefusedGroup:
		data := parseCallBackFewParam(prefixCallBackRefuseUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		res.Msg, err = h.userRefuseJoinGroup(userID, cbq.From.UserName, groupID)

	case isDeleteUserFromGroup:
		data := parseCallBackFewParam(prefixCallBackDelUserFromGr, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		res.EditMesage, err = h.deleteUserFromGroup(cbq.From.ID, userID, groupID, cbq.Message.MessageID)

	default:
		msg := h.createErrorMessaeg(cbq.From.ID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return nil
	}

	if err != nil {
		msg := h.createErrorMessaeg(cbq.From.ID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return err
	}
	h.response <- res
	return nil
}

func (h *Hub) isCommand(text string, msgInfo *tg.Message, timeStart time.Time) error {
	var msg *tg.MessageConfig
	var err error
	switch text {
	case "/start":
		msg, err = h.cmdStrart(msgInfo.From.ID, msgInfo.From.UserName)
		if err != nil {
			return err
		}

	case "/menu":
		msg = h.cmdGetMenu(msgInfo.From.ID)

	case "/help":
		msg = h.cmdHelp(msgInfo.Chat.ID)
	default:
		msg = h.cmdDefault(msg.ChatID)

	}
	if msg != nil {
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (h *Hub) isMessage(text string, msgInfo *tg.Message) (*tg.MessageConfig, error) {
	return nil, nil
}

func (h *Hub) isForwardMessage(msg *tg.Message, timeStart time.Time) error {
	text := msg.ReplyToMessage.Text
	var res *tg.MessageConfig
	var err error
	switch {
	case isCreateNameForward(text):

		res, err = h.createList(msg.Chat.ID, msg.Text)

	case isAddNewProductForward(text):
		listName := parseNameListActions(text)
		u := store.User{ChatID: msg.From.ID, UserName: &msg.From.UserName}
		res, err = h.addNewProduct(u, msg.Text, listName, false)

	case isAddNewProductGroupForward(text):
		listName := parseNameListActions(text)
		u := store.User{ChatID: msg.From.ID, UserName: &msg.From.UserName}
		res, err = h.addNewProduct(u, msg.Text, listName, true)

	case isEditListForward(text):
		listName := parseNameListActions(text)
		indexToDelete := parseIndexEditProduct(msg.Text)
		res, err = h.editProductList(msg.From.ID, listName, indexToDelete, false)

	case isEditGroupListForward(text):
		listName := parseNameListActions(text)
		indexToDelete := parseIndexEditProduct(msg.Text)
		res, err = h.editProductList(msg.From.ID, listName, indexToDelete, true)

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
	if res != nil {
		h.response <- MessageWithTime{Msg: res, WorkTime: timeStart}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (h *Hub) editReplyMarkup(chatID int64, markup *tg.InlineKeyboardMarkup, lastMsgID int) *tg.EditMessageReplyMarkupConfig {
	edits := tg.NewEditMessageReplyMarkup(chatID, lastMsgID, *markup)
	return &edits
}

func (h *Hub) editMessage(chatID int64, lastmsgDI int, text string) *tg.EditMessageTextConfig {
	msg := tg.NewEditMessageText(chatID, lastmsgDI, text)
	msg.ParseMode = "html"
	return &msg
}

func (h *Hub) createMessage(ChatId int64, text string) *tg.MessageConfig {
	msgCongig := tg.NewMessage(ChatId, text)
	msgCongig.ParseMode = "html"
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
