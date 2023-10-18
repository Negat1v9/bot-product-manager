package telegram

import (
	"context"
	"errors"
	"time"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

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
	textForMsg := createMessgeToInviteNewUser(ownerGroup.UserName, groupInfo.GroupName)
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
