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
	var msg *tg.MessageConfig
	var editMsg *tg.EditMessageTextConfig
	var err error
	CBType := h.getCallBackType(cbq.Data)
	switch CBType {

	case isGetLists:
		editMsg, err = h.getListName(cbq.From.ID, cbq.Message.MessageID)

	case isGetGroupProductList:
		data := parseCallBackFewParam(prefixCallBackGroupProductList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.getProductList(cbq.From.ID, cbq.Message.MessageID, listID, listName, true)

	case isGetGroupLists:
		editMsg, err = h.GetAllUserGroup(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateList:
		editMsg = h.createMsgToCreateList(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateGroup:
		editMsg, err = h.createMessageForNewGroup(cbq.From.ID, cbq.Message.MessageID)

	case isGetProductList:
		data := parseCallBackFewParam(prefixCallBackListProduct, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.getProductList(cbq.From.ID, cbq.Message.MessageID, listID, listName, false)

	case isWantConnectTemplate:
		data := parseCallBackFewParam(prefixWantConnectTemplate, cbq.Data)
		groupID, newListID := convSToI[int](data[0], 0), convSToI[int](data[1], 0)
		editMsg, err = h.getTemplatesForConnect(cbq.From.ID, groupID, newListID, cbq.Message.MessageID)

	case isGetTemplateForConnect:
		data := parseCallBackFewParam(prefixGetListForTemplateMerge, cbq.Data)
		listID, sNewID, sGrID := convSToI[int](data[0], 0), data[1], data[2]
		editMsg, err = h.getOneTemplateForConnect(cbq.From.ID, sNewID, sGrID, listID, cbq.Message.MessageID)

	case isConnectTemplate:
		data := parseCallBackFewParam(prefixConnectTemplate, cbq.Data)
		listID, newID, grID := convSToI[int](data[0], 0), convSToI[int](data[1], 0), convSToI[int](data[2], 0)
		editMsg, err = h.connectTemplate(cbq.From.ID, grID, listID, newID, cbq.Message.MessageID)

	case isWantAddNewProduct:
		products := cbq.Message.Text
		listName := parseCallBackOneParam(prefixAddProductList, cbq.Data)
		editMsg = h.wantAddNewProduct(cbq.From.ID, products, listName, cbq.Message.MessageID, false)
		// TODO:
	case isWantAddProductGroupList:
		products := cbq.Message.Text
		listName := parseCallBackOneParam(prefixAddProductGroup, cbq.Data)
		editMsg = h.wantAddNewProduct(cbq.From.ID, products, listName, cbq.Message.MessageID, true)

	case isGetAllGroupLists:
		data := parseCallBackOneParam(prefixCallBackListGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.GetGroupLists(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isGetGroupTemplates:
		data := parseCallBackOneParam(prefixGetGroupTemplates, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.getAllGroupTemplates(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isGetOneGroupTemplate:
		data := parseCallBackFewParam(prefixGetOneTemplate, cbq.Data)
		listID, listName, sGroupID := convSToI[int](data[0], 0), data[1], data[2]
		editMsg, err = h.getOneGroupTemplate(cbq.From.ID, sGroupID, listName, listID, cbq.Message.MessageID)

	case isWantEditList:
		products := cbq.Message.Text
		groupName := parseCallBackOneParam(prefixChangeList, cbq.Data)
		editMsg = h.createMessageForEditList(cbq.From.ID, products, groupName, cbq.Message.MessageID, false)
		// TODO:
	case isWantEditGroupList:
		products := cbq.Message.Text
		groupName := parseCallBackOneParam(prefixChangeGroupList, cbq.Data)
		editMsg = h.createMessageForEditList(cbq.From.ID, products, groupName, cbq.Message.MessageID, true)

	case isWantCreateGroupList:
		data := parseCallBackOneParam(prefixCreateGroupList, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.createMessageCreateGroupList(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isCompliteSoloList:
		data := parseCallBackFewParam(prefixCompliteSoloList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.compliteProductList(cbq.From.ID, listName, listID, cbq.Message.MessageID)

	case isWantCompliteList:
		products := cbq.Message.Text
		data := parseCallBackFewParam(prefixWantCompliteList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.wantCompliteList(cbq.From.ID, listName, products, listID, cbq.Message.MessageID)

	case isCompliteList:
		data := parseCallBackFewParam(prefixCompliteList, cbq.Data)
		listID, listName := convSToI[int](data[0], 0), data[1]
		editMsg, err = h.compliteProductList(cbq.From.ID, listName, listID, cbq.Message.MessageID)

	case isSaveTemplete:
		data := parseCallBackOneParam(prefixSaveAsTemplete, cbq.Data)
		listID := convSToI[int](data, 0)
		editMsg, err = h.saveAsTemplate(cbq.From.ID, listID, cbq.Message.MessageID)

	case isWantMergeList:
		listName := parseCallBackOneParam(prefixToMergeListGroup, cbq.Data)
		editMsg, err = h.getNameGroupMergeList(cbq.From.ID, listName, cbq.Message.MessageID)

	case isGetMainMenu:
		editMsg = h.getMainMenu(cbq.From.ID, cbq.Message.MessageID)

	case isMergeListGroup:
		data := parseCallBackFewParam(prefixMergeListWithGroup, cbq.Data)
		groupID, listID := convSToI[int](data[0], 0), convSToI[int](data[1], 0)
		editMsg, err = h.mergeListWithGroup(cbq.From.ID, groupID, listID, cbq.Message.MessageID)

	case isLeaveGroup:
		data := parseCallBackOneParam(prefixLeaveGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.leaveFromGroup(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isLeaveOwnerGroup:
		data := parseCallBackOneParam(prefixLeaveOwnerGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.leaveAndDeleteGroup(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isGetAllUsersGroup:
		data := parseCallBackOneParam(prefixGetAllUsersGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.getUserFromGroup(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isGetUsersToDelete:
		data := parseCallBackOneParam(prefixGetUserToDelete, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.getUserForDeleteFrGr(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isWantInviteNewUser:
		data := parseCallBackOneParam(prefixAddUserGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		editMsg, err = h.createMessageForInviteUser(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isUserReadyJoinGroup:
		data := parseCallBackFewParam(prefixCallBackInsertUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		msg, err = h.userReadyJoinGroup(userID, groupID)

	case isUserRefusedGroup:
		data := parseCallBackFewParam(prefixCallBackRefuseUserGroup, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		msg, err = h.userRefuseJoinGroup(userID, cbq.From.UserName, groupID)

	case isDeleteUserFromGroup:
		data := parseCallBackFewParam(prefixCallBackDelUserFromGr, cbq.Data)
		userID, groupID := convSToI[int64](data[0], 64), convSToI[int](data[1], 0)
		editMsg, err = h.deleteUserFromGroup(cbq.From.ID, userID, groupID, cbq.Message.MessageID)

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
	h.response <- MessageWithTime{Msg: msg, EditMesage: editMsg, WorkTime: timeStart}
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
		list := &store.ProductList{
			OwnerID: &msg.From.ID,
			Name:    &msg.Text,
		}
		res, err = h.createList(msg.Chat.ID, list)

	case isAddNewProductForward(text):
		listName := parseNameListActions(text)

		res, err = h.addNewProduct(msg.Chat.ID, msg.Text, listName, false)

	case isAddNewProductGroupForward(text):
		listName := parseNameListActions(text)
		res, err = h.addNewProduct(msg.Chat.ID, msg.Text, listName, true)

	case isEditListForward(text):
		listName := parseNameListActions(text)
		indexToDelete := parseIndexEditProduct(msg.Text)
		res, err = h.editProductList(msg.From.ID, listName, indexToDelete, false)

	case isEditGroupListForward(text):
		listName := parseNameListActions(text)
		indexToDelete := parseIndexEditProduct(msg.Text)
		res, err = h.editProductList(msg.From.ID, listName, indexToDelete, false)

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
