package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: create message with information to user, about he have to forwad its
func (h *Hub) answerToCreateList(ChatID int64) *tg.MessageConfig {
	msg := h.createMessage(ChatID, answerCreateListMsg)
	return msg
}

// Info: Create list in database, and answer with new list name
func (h *Hub) createList(ChatID int64, list *store.ProductList) (*tg.MessageConfig, error) {

	err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", *list.Name))
	return msg, nil
}

// Info: Select all users lists and create inline keyboard with its
func (h *Hub) getListName(UserID, ChatID int64) (msg *tg.MessageConfig, err error) {
	lists, err := h.db.ProductList().GetAll(context.TODO(), int(UserID))
	if err != nil {
		if err == store.NoRowListOfProductError {
			msg = h.createMessage(ChatID, "Nothing is found. Create Youre First list!")
			return msg, nil
		}
		return nil, err
	}
	keyboard := createListProductInline(lists)
	msg = h.createMessage(ChatID, "List of Product-lists")
	msg.ReplyMarkup = keyboard
	return msg, nil
}

func (h *Hub) getProductList(ChatID int64, listID int, listName string) (msg *tg.MessageConfig, err error) {
	product, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		if err == store.NoRowProductError {

			msg = h.createMessage(ChatID, emptyListMessage)
			msg.ReplyMarkup = createProductsInline(listName)
			return msg, nil

		}
		return nil, err
	}
	text := createMessageProductList(product.Products)
	msg = h.createMessage(ChatID, text)

	msg.ReplyMarkup = createProductsInline(listName)
	return msg, nil
}

func (h *Hub) addNewProduct(ChatID int64, p store.Product, listName string) (*tg.MessageConfig, error) {
	product, err := h.db.Product().GetAll(context.TODO(), p.ListID)
	if err != nil {
		// if not row exist
		if err == store.NoRowProductError {
			err = h.db.Product().Create(context.TODO(), p.ListID)
			product = &store.Product{
				ListID:   p.ListID,
				Products: []string{},
			}
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	text := createMessageSuccessAddedProduct(p.Products)
	// add to old values new
	if len(product.Products) > 0 {
		product.Products = append(product.Products, p.Products...)
		p.Products = product.Products
	}
	if err := h.db.Product().Add(context.TODO(), p); err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, text)
	msg.ReplyMarkup = createInlineGetCurList(p.ListID, listName)
	return msg, nil
}

func (h *Hub) compliteProductList(ChatID int64, productListName string) (*tg.MessageConfig, error) {
	listID, err := h.db.ProductList().GetListID(context.TODO(), productListName)
	if err != nil {
		return nil, err
	}
	err = h.db.ProductList().Delete(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, isCompletesProductListMsg+productListName)
	return msg, nil
}

func (h *Hub) createMessageForNewGroup(ChatID int64) *tg.MessageConfig {
	msg := h.createMessage(ChatID, createGroupMessage)
	return msg
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
	return msg, nil
}

func (h *Hub) GetAllUserGroup(ChatID, UserID int64) (*tg.MessageConfig, error) {
	groups, err := h.db.ManagerGroup().UserGroup(context.TODO(), int(UserID))
	if err != nil {
		if err == store.NoUserGroupError {
			msg := h.createMessage(ChatID, err.Error())
			return msg, nil
		}
		return nil, err
	}
	msg := h.createMessage(ChatID, `Groups:`)
	msg.ReplyMarkup = createInlineGroupName(groups)
	return msg, nil
}

func (h *Hub) GetGroupLists(UserID int64, groupID int) (msg *tg.MessageConfig, err error) {
	var isOwnerGroup = false
	groupList, err := h.db.ManagerGroup().AllByGroupID(context.TODO(), groupID)
	if err != nil {
		if err == store.NoRowListOfProductError {
			msg = h.createMessage(UserID, err.Error())
		} else {
			return nil, err
		}
	} else {
		msg = h.createMessage(UserID, "group lists")
	}
	if UserID == groupList.GroupOwnerID {
		isOwnerGroup = true
	}
	msg.ReplyMarkup = createInlineGroupList(groupList.PruductLists, groupID, isOwnerGroup)
	return msg, nil
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

func (h *Hub) getUserForDeleteFrGr(ChatID int64, groupID int) (msg *tg.MessageConfig, err error) {
	groupInfo, err := h.db.ManagerGroup().InfoGroup(context.TODO(), groupID)
	if err != nil {
		return nil, err
	}
	// Info: if user at group lowwer them 2, means at group only owner
	if len(*groupInfo.UsersInfo) < 2 {
		msg = h.createMessage(ChatID, emptyUserInGroup)
		return
	}
	msg = h.createMessage(ChatID, "Choise user from:"+groupInfo.GroupName)
	msg.ReplyMarkup = createInlineDeleteUser(*groupInfo.UsersInfo, groupID)
	return
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
