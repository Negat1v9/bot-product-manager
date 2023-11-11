package telegram

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: Send edit message with all current group lists
func (h *Hub) GetGroupLists(UserID int64, lastMsgID, groupID int) (editMsg *tg.EditMessageTextConfig, err error) {
	groupList, err := h.db.ManagerGroup().AllByGroupID(context.TODO(), groupID)
	if err != nil {
		if err == store.NoRowListOfProductError {
			editMsg = h.editMessage(UserID, lastMsgID, err.Error())
		} else {
			return nil, err
		}
	} else {
		editMsg = h.editMessage(UserID, lastMsgID, "Group List:")
	}
	// isOwnerGroup := UserID == groupList.GroupOwnerID

	editMsg.ReplyMarkup = createInlineGroupList(groupList.PruductLists, groupID)
	return editMsg, nil
}

// Info: Create message with info for reply message to create new group list
func (h *Hub) createMessageCreateGroupList(UserID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(UserID, lastMsgID, answerCreateGroupListMsg+group.GroupName)
	return editMsg, nil
}

// Info: Create group list in database
func (h *Hub) createGroupList(UserID int64, listName, groupName string) (*tg.MessageConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupName(context.TODO(), groupName)
	if err != nil {
		return nil, err
	}
	list := &store.ProductList{
		OwnerID: &UserID,
		GroupID: &group.ID,
		Name:    &listName,
	}
	err = h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(UserID, `New list is created`)
	msg.ReplyMarkup = createInlineGetCurGroup(group.ID)
	return msg, nil
}

func (h *Hub) getUserFromGroup(chatID int64, lastMsgID, groupID int) (*tg.EditMessageTextConfig, error) {
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	text := createMessageGetAllUsersGroup(*groupInfo.UsersInfo, groupInfo.OwnerID)
	editMsg := h.editMessage(chatID, lastMsgID, text)
	isOwnerGroup := chatID == groupInfo.OwnerID
	// if user is owner he can invite and delete new and old users
	if isOwnerGroup {
		editMsg.ReplyMarkup = creaetInlineUsersGroupActions(groupInfo.ID)
	} else {
		editMsg.ReplyMarkup = creaetInlineBackToGroupButton(groupID)
	}
	return editMsg, nil
}

// Info: Get list with all users in group with delete button
func (h *Hub) getUserForDeleteFrGr(ChatID int64, lastMsgID, groupID int) (*tg.EditMessageTextConfig, error) {
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	if len(*groupInfo.UsersInfo) < 2 {
		editMsg := h.editMessage(ChatID, lastMsgID, emptyUserInGroup)
		editMsg.ReplyMarkup = creaetInlineBackToGroupButton(groupID)
		return editMsg, nil
	}
	editMsg := h.editMessage(ChatID, lastMsgID, "Choise user from:"+groupInfo.GroupName)
	editMsg.ReplyMarkup = createInlineDeleteUser(*groupInfo.UsersInfo, groupID, groupInfo.OwnerID)
	return editMsg, nil
}

// Info: User push on botton with name user who will deleted
func (h *Hub) deleteUserFromGroup(ChatID, userID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	g := &store.Group{
		UserID:  userID,
		GroupID: groupID,
	}
	err := h.db.Group().DeleteUser(context.TODO(), g)
	if err != nil {
		return nil, err
	}
	msg := h.editMessage(ChatID, lastMsgID, successDeletedUser)
	return msg, nil
}

// Info: create message to reply on this, manager will invite new user by nickname
func (h *Hub) createMessageForInviteUser(chatID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, textForInvitingNewUser+group.GroupName)
	return editMsg, nil
}
