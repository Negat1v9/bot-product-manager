package telegram

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Hub) createMessageForNewGroup(ChatID int64, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	var editMsg *tg.EditMessageTextConfig
	user, err := h.db.User().ByID(context.TODO(), ChatID)
	if err != nil {
		return nil, err
	}
	if user.UserName != nil {
		editMsg = h.editMessage(ChatID, lastMsgID, createGroupMessage)
	} else {
		editMsg = h.editMessage(ChatID, lastMsgID, NotNickNameUserMsg)
	}
	return editMsg, nil
}

func (h *Hub) createNewGroup(ChatID int64, managerGroup *store.GroupInfo) (*tg.MessageConfig, error) {
	id, err := h.db.ManagerGroup().Create(context.TODO(), managerGroup)
	if err != nil {
		return nil, err
	}
	group := &store.Group{
		UserID:  managerGroup.OwnerID,
		GroupID: id,
	}
	err = h.db.Group().AddUser(context.TODO(), group)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, groupIsCreatesMessage)
	msg.ReplyMarkup = createInlineGoToGroups()
	return msg, nil
}

func (h *Hub) GetAllUserGroup(chatID int64, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	groups, err := h.db.ManagerGroup().UserGroup(context.TODO(), chatID)
	if err != nil {
		if err == store.NoUserGroupError {
			editMsg := h.editMessage(chatID, lastMsgID, err.Error())
			editMsg.ReplyMarkup = createInlineGoToMenu()
			return editMsg, nil
		}
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "All the groups you are in 🎡")
	editMsg.ReplyMarkup = createInlineGroupName(groups)
	return editMsg, nil
}

func (h *Hub) inviteNewUser(ChatID int64, newUserName, groupName string) (*tg.MessageConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupName(context.TODO(), groupName)
	if err != nil {
		return nil, err
	}
	err = h.sendInviteMessage(newUserName, group.ID, ChatID)
	if err != nil {
		if err == userAlredyGroupError {
			msg := h.createMessage(ChatID, err.Error())
			return msg, nil
		}
		// send invited message for forward
		if err == userNoExistError {
			msg := h.createMessage(ChatID, joinNewUserAtBotMessage)
			return msg, nil
		}
	}
	msg := h.createMessage(ChatID, inviteSendMessage)
	return msg, nil
}

func (h *Hub) userReadyJoinGroup(newUserID int64, groupID int) (*tg.MessageConfig, error) {
	g := &store.Group{
		UserID:  newUserID,
		GroupID: groupID,
	}
	err := h.db.Group().AddUser(context.TODO(), g)
	if err != nil {
		return nil, err
	}
	editMsg := h.createMessage(newUserID, userInvitedInGroupMessage)
	return editMsg, nil
}

func (h *Hub) userRefuseJoinGroup(userID int64, userName string, groupID int) (*tg.MessageConfig, error) {
	if err := h.sendRufuseAnswerToOwner(groupID, userName); err != nil {
		return nil, err
	}
	editMsg := h.createMessage(userID, refuseJoinGroupMessage)
	return editMsg, nil
}

func (h *Hub) leaveFromGroup(chatID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.OwnerID == chatID {
		editMsg := h.editMessage(chatID, lastMsgID, ownerGroupWantLeave)
		editMsg.ReplyMarkup = createInlineMakeSureDelete(groupID)
		return editMsg, nil
	}
	g := &store.Group{
		UserID:  chatID,
		GroupID: groupInfo.ID,
	}
	// delete user
	if err = h.db.Group().DeleteUser(context.TODO(), g); err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, successLeaveGroup)
	editMsg.ReplyMarkup = createInlineGoToGroups()
	return editMsg, nil
}

func (h *Hub) leaveAndDeleteGroup(chatID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	err := h.db.ManagerGroup().DeleteGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, successLeaveGroup)
	editMsg.ReplyMarkup = createInlineGoToGroups()
	return editMsg, nil
}
