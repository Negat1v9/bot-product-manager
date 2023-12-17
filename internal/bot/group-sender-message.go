package telegram

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

const timeOutForComliteList = 20

var (
	userNoExistError     = errors.New("user not exist")
	userAlredyGroupError = errors.New("user is alredy in group")
)

func (h *Hub) sendInviteMessage(newUserName string, groupID int, ownerID int64) error {
	newUser, err := h.db.User().ByUserName(context.TODO(), newUserName)
	if err != nil {
		if err == store.NoExistUserError {
			return userNoExistError
		}
		return err
	}

	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return err
	}

	if checkUserInGroup(newUser.ChatID, *groupInfo.UsersInfo) {
		return userAlredyGroupError
	}
	ownerGroup := searchOwnerGroup(ownerID, *groupInfo.UsersInfo)
	textForMsg := createMessgeToInviteNewUser(*ownerGroup.UserName, groupInfo.GroupName)
	msgForNewUser := h.createMessage(newUser.ChatID, textForMsg)
	msgForNewUser.ReplyMarkup = createInlineInviteUserGroup(groupID, newUser.ChatID)
	h.response <- MessageWithTime{Msg: msgForNewUser, WorkTime: time.Now()}
	return nil
}

func (h *Hub) sendRufuseAnswerToOwner(groupID int, refusedName string) error {
	start := time.Now()
	groupInfo, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return err
	}
	ownerUser, err := h.db.User().ByID(context.TODO(), groupInfo.OwnerID)
	if err != nil {
		return nil
	}
	text := createMessageUserRefusedOrder(refusedName)
	msg := h.createMessage(ownerUser.ChatID, text)
	h.response <- MessageWithTime{Msg: msg, WorkTime: start}
	return nil
}

func (h *Hub) sendNotifAddNewList(createrID int64, groupID int, listName string) error {
	start := time.Now()
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return err
	}
	for _, u := range *groupInfo.UsersInfo {
		if u.ChatID == createrID {
			continue
		}
		text := fmt.Sprintf(
			"⚡ Hello <b>%s</b>, a new list has appeared in the %s group.\nLook what it is 👀",
			*u.UserName,
			groupInfo.GroupName,
		)
		msg := h.createMessage(u.ChatID, text)
		msg.ReplyMarkup = createInlineGetCurGroup(groupID)
		h.response <- MessageWithTime{Msg: msg, WorkTime: start}
	}
	return nil
}

func (h *Hub) sendComplitedListGroupDelay(userComlite int64, listID, groupID int) {
	if ok := h.setTimerForComliteMsg(listID); ok {
		return
	}
	timeStart := time.Now()
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	// FIXME: What if error?
	if err != nil {
		return
	}
	infoList, err := h.db.ProductList().GetFoolInfoGroupProdList(context.TODO(), listID)
	if err != nil {
		return
	}
	text := createMessageCompliteGroupList(*infoList, userComlite)
	if err = h.db.ProductList().Delete(context.TODO(), listID); err != nil {
		return
	}
	for _, user := range *groupInfo.UsersInfo {
		if userComlite == user.ChatID {
			continue
		}
		msg := h.createMessage(user.ChatID, text)
		h.response <- MessageWithTime{Msg: msg, WorkTime: timeStart}
	}

}

// Info: get user time for recover list
// true - the list is was recovered
// false - user confirm list is ready
func (h *Hub) setTimerForComliteMsg(listID int) bool {

	h.container.SetRecoverList(listID)

	defer h.container.DeleteRecoverList(listID)

	time.Sleep(time.Second * timeOutForComliteList)

	if h.container.isInContainerList(listID) {
		return false
	}
	return true
}
