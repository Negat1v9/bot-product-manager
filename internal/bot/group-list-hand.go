package telegram

import (
	"context"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Hub) GetGroupLists(UserID int64, lastMsgID, groupID int) (editMsg *tg.EditMessageTextConfig, err error) {
	var isOwnerGroup = false
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
	if UserID == groupList.GroupOwnerID {
		isOwnerGroup = true
	}
	editMsg.ReplyMarkup = createInlineGroupList(groupList.PruductLists, groupID, isOwnerGroup)
	return editMsg, nil
}

func (h *Hub) createMessageCreateGroupList(UserID int64, groupID int) (*tg.MessageConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(UserID, answerCreateGroupListMsg+group.GroupName)
	return msg, nil
}

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
	return msg, nil
}

func (h *Hub) getUserForDeleteFrGr(ChatID int64, lastMsgID, groupID int) (*tg.EditMessageTextConfig, error) {
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	// Info: if user at group lowwer them 2, means at group only owner
	if len(*groupInfo.UsersInfo) < 2 {
		editMsg := h.editMessage(ChatID, lastMsgID, emptyUserInGroup)
		// msg = h.createMessage(ChatID, emptyUserInGroup)
		return editMsg, nil
	}
	editMsg := h.editMessage(ChatID, lastMsgID, "Choise user from:"+groupInfo.GroupName)
	// msg = h.createMessage(ChatID, "Choise user from:"+groupInfo.GroupName)
	editMsg.ReplyMarkup = createInlineDeleteUser(*groupInfo.UsersInfo, groupID)
	return editMsg, nil
}

func (h *Hub) deleteUserFromGroup(ChatID, userID int64, groupID int) (*tg.MessageConfig, error) {
	g := &store.Group{
		UserID:  userID,
		GroupID: groupID,
	}
	err := h.db.Group().DeleteUser(context.TODO(), g)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, successDeletedUser)
	return msg, nil
}

func (h *Hub) createMessageForInviteUser(chatID int64, groupID int) (*tg.MessageConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(chatID, textForInvitingNewUser+group.GroupName)
	return msg, nil
}
