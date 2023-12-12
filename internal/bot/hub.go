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

func (h *Hub) MessageUpdate(msg *tg.Message, timeStart time.Time) (err error) {
	text := toLowerCase(msg.Text)

	typeUpdate := h.getTypeMessage(text, msg)

	switch typeUpdate {

	case CommandUpdate:
		err = h.isCommand(text, msg, timeStart)

	case TextUpdate:

		t, ok := h.getTypeUserCmd(msg.From.ID)
		if !ok {
			msg := h.createMessage(msg.From.ID, errorLastCmdUserDeleted)
			h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
			return
		}
		err = h.isMessageText(msg, timeStart, *t)
	}
	if err != nil {

		editMsg := h.createErrorMessaeg(msg.From.ID, msg.MessageID)
		h.response <- MessageWithTime{EditMesage: editMsg, WorkTime: timeStart}
		return err
	}
	return nil
}

func (h *Hub) getTypeMessage(text string, msg *tg.Message) int {
	if isCommandUpdate(text) {
		return CommandUpdate
	}

	return TextUpdate
}

func (h *Hub) getTypeUserCmd(userID int64) (*TypeUserCommand, bool) {
	t, ok := h.container.UsersCmd[userID]
	if ok {
		h.container.DeleteUserCmd(userID)
		return &t, ok
	}
	return nil, ok
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
		data := parseCallBackOneParam(prefixCallBackGroupProductList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage, err = h.getGroupList(cbq.From.ID, cbq.Message.MessageID, listID)

	case isGetGroupLists:
		res.EditMesage, err = h.GetAllUserGroup(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateList:
		res.EditMesage = h.createMsgToCreateList(cbq.From.ID, cbq.Message.MessageID)

	case isWantCreateGroup:
		res.EditMesage, err = h.createMessageForNewGroup(cbq.From.ID, cbq.Message.MessageID)

	case isGetProductList:
		data := parseCallBackOneParam(prefixCallBackListProduct, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage, err = h.getProductListV2(cbq.From.ID, cbq.Message.MessageID, listID)

	case isWantAddNewProduct:
		products := cbq.Message.Text
		data := parseCallBackOneParam(prefixAddProductList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage = h.wantAddNewProduct(cbq.From.ID, products, listID, cbq.Message.MessageID, false)

	case isWantAddProductGroupList:
		products := cbq.Message.Text
		data := parseCallBackOneParam(prefixAddProductGroup, cbq.Data)
		listID := convSToI[int](data, 0)

		res.EditMesage = h.wantAddNewProduct(cbq.From.ID, products, listID, cbq.Message.MessageID, true)

	case isGetAllGroupLists:
		data := parseCallBackOneParam(prefixCallBackListGroup, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.GetGroupLists(cbq.From.ID, cbq.Message.MessageID, groupID)

	case isWantDeleteProductFromList:
		data := parseCallBackFewParam(prefixGetPageProdDelete, cbq.Data)
		listID, offset := convSToI[int](data[0], 0), convSToI[int](data[1], 0)
		res.EditMesage, err = h.wantDeleteProd(cbq.From.ID, listID, offset, cbq.Message.MessageID, false)

	case isDeleteProd:
		data := parseCallBackFewParam(prefixDeleteProd, cbq.Data)
		prodID, sListID := convSToI[int](data[0], 0), data[1]
		res.EditMesage, err = h.deleteProd(cbq.From.ID, prodID, cbq.Message.MessageID, sListID, false)

	case isWantDeleteProductFromGroupList:

		data := parseCallBackFewParam(prefixGetPageGroupProdDelete, cbq.Data)
		listID, offset := convSToI[int](data[0], 0), convSToI[int](data[1], 0)
		res.EditMesage, err = h.wantDeleteProd(cbq.From.ID, listID, offset, cbq.Message.MessageID, true)

	case isDeleteGroupProd:
		data := parseCallBackFewParam(prefixDeleteGrProd, cbq.Data)
		prodID, sListID := convSToI[int](data[0], 0), data[1]
		res.EditMesage, err = h.deleteProd(cbq.From.ID, prodID, cbq.Message.MessageID, sListID, true)

	case isWantCreateGroupList:
		data := parseCallBackOneParam(prefixCreateGroupList, cbq.Data)
		groupID := convSToI[int](data, 0)
		res.EditMesage, err = h.createMessageCreateGroupList(cbq.From.ID, groupID, cbq.Message.MessageID)

	case isCompliteSoloList:
		data := parseCallBackOneParam(prefixCompliteSoloList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage, err = h.compliteProductList(cbq.From.ID, listID, cbq.Message.MessageID)

	case isWantCompliteGrList:
		products := cbq.Message.Text
		data := parseCallBackOneParam(prefixWantCompliteGrList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage, err = h.wantCompliteList(cbq.From.ID, products, listID, cbq.Message.MessageID)

	case isCompliteGroupList:
		data := parseCallBackFewParam(prefixCompliteList, cbq.Data)
		listID, sGroupID, listName := convSToI[int](data[0], 0), data[1], data[2]
		res.EditReplyMarkup, err = h.compliteGroupList(cbq.From.ID, listName, sGroupID, cbq.Message.Text, listID, cbq.Message.MessageID)

	case isRestoreProductList:
		data := parseCallBackOneParam(prefixRestoreList, cbq.Data)
		listID := convSToI[int](data, 0)
		res.Msg, err = h.recoverUserList(cbq.From.ID, listID, cbq.Message.Text)

	case isRestoreGroupList:
		data := parseCallBackFewParam(prefixRestoreGroupList, cbq.Data)
		listID, groupID, listName := convSToI[int](data[0], 0), convSToI[int](data[1], 0), data[2]

		res.Msg, err = h.recoverGroupList(cbq.From.ID, listID, groupID, cbq.Message.MessageID, cbq.Message.Text, listName)
	case isWantMergeList:
		data := parseCallBackOneParam(prefixToMergeListGroup, cbq.Data)
		listID := convSToI[int](data, 0)
		res.EditMesage, err = h.getNameGroupMergeList(cbq.From.ID, listID, cbq.Message.MessageID)

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
		res.EditMesage = h.createErrorMessaeg(cbq.From.ID, cbq.Message.MessageID)
		h.response <- res
		return nil
	}

	if err != nil {
		res.EditMesage = h.createErrorMessaeg(cbq.From.ID, cbq.Message.MessageID)
		h.response <- res
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

func (h *Hub) isMessageText(msg *tg.Message, timeStart time.Time, typeCmd TypeUserCommand) error {
	var res *tg.MessageConfig
	var err error
	switch typeCmd.TypeCmd {

	case isCreateNewList:
		res, err = h.createList(msg.Chat.ID, msg.Text)

	case isAddNewProduct:

		u := store.User{ChatID: msg.From.ID, UserName: &msg.From.UserName}
		res, err = h.addNewProductV2(u, msg.Text, *typeCmd.ListID, false)

	case isAddNewProductGroup:

		u := store.User{ChatID: msg.From.ID, UserName: &msg.From.UserName}
		res, err = h.addNewProductV2(u, msg.Text, *typeCmd.ListID, true)

	case isCreateGroup:
		managerGroup := &store.GroupInfo{
			OwnerID:   msg.From.ID,
			GroupName: msg.Text,
		}
		res, err = h.createNewGroup(msg.Chat.ID, managerGroup)

	case isCreateGroupList:

		res, err = h.createGroupList(msg.From.ID, msg.Text, *typeCmd.GroupID)

	case isSendInviteNewUser:
		newUserName := parseUserNickNameForAddGroup(msg.Text)
		res, err = h.inviteNewUser(msg.From.ID, newUserName, *typeCmd.GroupID)
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

func (h *Hub) createErrorMessaeg(chatID int64, lastMsgID int) *tg.EditMessageTextConfig {
	editMsg := h.editMessage(chatID, lastMsgID, errorMessage)
	return editMsg
}

// NOTE: Make my example
func toLowerCase(s string) string {
	return strings.ToLower(s)
}
