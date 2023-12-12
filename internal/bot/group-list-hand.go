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

func (h *Hub) getGroupList(chatID int64, lastMsgID, listID int) (em *tg.EditMessageTextConfig, err error) {
	prod, err := h.db.Product().GetByListID(context.TODO(), listID, 0, 100)
	if err != nil {
		if err == store.NoRowProductError {
			em = h.editMessage(chatID, lastMsgID, emptyListMessage)

		} else {
			return nil, err
		}

	} else if len(prod) == 0 {
		em = h.editMessage(chatID, lastMsgID, emptyListMessage)
	} else {
		text := createMessageProductList(prod)
		em = h.editMessage(chatID, lastMsgID, text)
	}
	em.ReplyMarkup = createInlineProductsGroup(listID)
	return em, nil
}

// Info: Create message with info for reply message to create new group list
func (h *Hub) createMessageCreateGroupList(chatID int64, groupID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, answerCreateGroupListMsg+group.GroupName)
	userCmd := TypeUserCommand{
		TypeCmd: isCreateGroupList,
		GroupID: &groupID,
	}
	h.container.AddUserCmd(chatID, userCmd)
	return editMsg, nil
}

// Info: Create group list in database
func (h *Hub) createGroupList(UserID int64, listName string, groupID int) (*tg.MessageConfig, error) {
	group, err := h.db.ManagerGroup().ByGroupID(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	clearNameList := makeNameClear(listName)
	list := &store.ProductList{
		OwnerID: &UserID,
		GroupID: &groupID,
		Name:    &clearNameList,
	}
	id, err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	go h.sendNotifAddNewList(UserID, groupID, group.GroupName)
	msg := h.createMessage(UserID, clearNameList+" created successfully ✅")

	msg.ReplyMarkup = createInlineGetCurGroupList(id)
	return msg, nil
}

func (h *Hub) wantCompliteList(chatID int64, products string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	if products == emptyListMessage {
		err := h.db.ProductList().Delete(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
		editMsg := h.editMessage(chatID, lastMsgID, isCompletesProductListMsg)
		editMsg.ReplyMarkup = createInlineGoToGroups()
		if err != nil {
			return nil, err
		}
		return editMsg, nil
	}
	infoList, err := h.db.ProductList().GetFoolInfoGroupProdList(context.TODO(), listID)
	// list, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		return nil, err
	}

	text := createMessageCompliteGroupList(*infoList, chatID)
	msg := h.editMessage(chatID, lastMsgID, text)

	msg.ReplyMarkup = createInlineAfterComplite(listID, *infoList.List.GroupID, *infoList.List.Name)
	return msg, nil
}

func (h *Hub) compliteGroupList(ChatID int64, name, sGrID, text string, listID, lastMsgID int) (*tg.EditMessageReplyMarkupConfig, error) {

	err := h.db.ProductList().MakeListInactive(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	var markup *tg.InlineKeyboardMarkup

	groupID := convSToI[int](sGrID, 0)
	go h.sendComplitedListGroupDelay(ChatID, listID, groupID, text)
	markup = createInlineRecoverGroupList(listID, sGrID, name)

	msg := h.editReplyMarkup(ChatID, markup, lastMsgID)
	return msg, nil
}

func (h *Hub) recoverGroupList(chatID int64, listID, groupID, lastMsgID int, text, listName string) (msg *tg.MessageConfig, err error) {
	if h.container.isInContainerList(listID) {

		h.container.DeleteRecoverList(listID)

		listName, err = h.db.ProductList().MakeListActive(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
	} else {
		list := &store.ProductList{
			Name:    &listName,
			OwnerID: &chatID,
			GroupID: &groupID,
		}
		listID, err := h.db.ProductList().Create(context.TODO(), list)
		if err != nil {
			return nil, err
		}
		splitedText := splitText(text, '\n')
		prod := parseTextToProd(splitedText, chatID, listID)
		err = h.db.Product().Add(prod)
		if err != nil {
			return nil, err
		}

	}
	msg = h.createMessage(chatID, "The list - "+listName+" is recover")
	msg.ReplyMarkup = createInlineGetCurGroupList(listID)
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
	editMsg := h.editMessage(ChatID, lastMsgID, "Choise user from: "+groupInfo.GroupName)
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
	userCmd := TypeUserCommand{
		TypeCmd: isSendInviteNewUser,
		GroupID: &groupID,
	}
	h.container.AddUserCmd(chatID, userCmd)
	return editMsg, nil
}
