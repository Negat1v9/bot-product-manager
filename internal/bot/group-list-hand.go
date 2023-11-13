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
			editMsg = h.editMessage(UserID, lastMsgID, "It looks like there are no lists in this group 😱, you can create the first one 💢")
		} else {
			return nil, err
		}
	} else {
		editMsg = h.editMessage(UserID, lastMsgID, "These are your group's lists 👇")
	}

	editMsg.ReplyMarkup = createInlineGroupList(groupList.PruductLists, groupID)
	return editMsg, nil
}

func (h *Hub) getAllGroupTemplates(chatID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	lists, err := h.db.ProductList().GetAllGroupTemplates(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "select a template to edit it 📄")
	editMsg.ReplyMarkup = createInlineGroupListTemplates(lists, groupID)
	return editMsg, nil
}

func (h *Hub) getOneGroupTemplate(chatID int64, sGroupID, listName string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	products, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		if err == store.NoRowProductError {
			// If list is empty notify user about this
			editMsg := h.editMessage(chatID, lastMsgID, emptyListMessage)
			editMsg.ReplyMarkup = createInlineTemplateActions(listID, listName, sGroupID)
			return editMsg, nil

		}
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, createMessageProductList(products.Products))
	editMsg.ReplyMarkup = createInlineTemplateActions(listID, listName, sGroupID)
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
	id, err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	go h.sendNotifAddNewList(UserID, group.ID, listName)
	msg := h.createMessage(UserID, getInformationMergeTemplateMsg)
	msg.ReplyMarkup = createInlineAfterListCreated(group.ID, id)
	return msg, nil
}

func (h *Hub) getTemplatesForConnect(chatID int64, groupID, newListID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	lists, err := h.db.ProductList().GetAllGroupTemplates(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "select template for connect ♻")
	editMsg.ReplyMarkup = createInlineTemplateForConnect(lists, groupID, newListID)
	return editMsg, nil
}

func (h *Hub) getOneTemplateForConnect(chatID int64, sNewID, sGrID string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	products, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	text := createMessageProductList(products.Products)
	editMsg := h.editMessage(chatID, lastMsgID, text)
	editMsg.ReplyMarkup = createInlineConnectTemplate(listID, sNewID, sGrID)
	return editMsg, nil
}

func (h *Hub) connectTemplate(chatID int64, groupID, listID, newID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	templateList, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	var newList *store.Product
	newList, err = h.db.Product().GetAll(context.TODO(), newID)
	// if no exist products in products table its means the list is clear
	if err != nil {
		if err == store.NoRowProductError {
			err := h.db.Product().Create(context.TODO(), newID)
			if err != nil {
				return nil, err
			}
			newList = &store.Product{
				ListID:   newID,
				Products: []string{},
			}
		} else {
			return nil, err
		}
	}
	// add products from template at list
	newList.Products = append(newList.Products, templateList.Products...)
	err = h.db.Product().Add(context.TODO(), *newList)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "✅ products from the template have been <b>successfully added</b> to the new list")
	editMsg.ReplyMarkup = createInlineGetCurGroup(groupID)
	return editMsg, nil
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
