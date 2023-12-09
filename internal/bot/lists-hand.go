package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Hub) createMsgToCreateList(ChatID int64, lastMsgID int) *tg.EditMessageTextConfig {
	typeUserCmd := TypeUserCommand{
		TypeCmd: isCreateNewList,
	}
	h.container.AddUserCmd(ChatID, typeUserCmd)
	editMsg := h.editMessage(ChatID, lastMsgID, answerCreateListMsg)
	return editMsg
}

// Info: Create list in database, and answer with new list name
func (h *Hub) createList(ChatID int64, nameList string) (*tg.MessageConfig, error) {
	clearName := makeNameClear(nameList)
	list := &store.ProductList{
		OwnerID: &ChatID,
		Name:    &clearName,
	}
	_, err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", *list.Name))
	msg.ReplyMarkup = createInlineGoToMenu()
	return msg, nil
}

// Info: Select all users lists and create inline keyboard with its
func (h *Hub) getListName(chatID int64, lastMsgID int) (editMsg *tg.EditMessageTextConfig, err error) {
	lists, err := h.db.ProductList().GetAllNames(context.TODO(), chatID)
	if err != nil {
		if err == store.NoRowListOfProductError {
			editMsg = h.editMessage(chatID, lastMsgID, "Nothing is found. Create Youre First list!")
			editMsg.ReplyMarkup = createInlineGoToMenu()
			return editMsg, nil
		}
		return nil, err
	}
	keyboard := createListProductInline(lists)
	editMsg = h.editMessage(chatID, lastMsgID, listsProductsMsgHelp)

	editMsg.ReplyMarkup = &keyboard
	return editMsg, nil
}

func (h *Hub) getProductList(ChatID int64, lastMsgID, listID int, listName string) (editMsg *tg.EditMessageTextConfig, err error) {
	productList, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		if err == store.NoRowProductError {
			editMsg = h.editMessage(ChatID, lastMsgID, emptyListMessage)

		} else {
			return nil, err
		}

	} else if len(productList.Products) == 0 {
		editMsg = h.editMessage(ChatID, lastMsgID, emptyListMessage)
	} else {
		text := createMessageProductList(productList.Products)
		editMsg = h.editMessage(ChatID, lastMsgID, text)
	}

	editMsg.ReplyMarkup = createProductsInline(listName, listID)
	return editMsg, nil
}

func (h *Hub) wantAddNewProduct(chatID int64, products, listName string, listID, lastMsgID int, isGroup bool) *tg.EditMessageTextConfig {
	var text string
	if products == emptyListMessage {
		text = addNewProductMessageReply + listName
	} else {
		text = products + "\n❓\n" + addNewProductMessageReply + listName
	}
	var userCmd TypeUserCommand
	if isGroup {
		userCmd = TypeUserCommand{
			TypeCmd: isAddNewProductGroup,
			ListID:  &listID,
		}
	} else {
		userCmd = TypeUserCommand{
			TypeCmd: isAddNewProduct,
			ListID:  &listID,
		}
	}
	h.container.AddUserCmd(chatID, userCmd)
	editMsg := h.editMessage(chatID, lastMsgID, text)
	return editMsg
}

func (h *Hub) addNewProduct(u store.User, products string, listID int, isGroup bool) (*tg.MessageConfig, error) {
	list, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		return nil, err
	}
	newProduct := parseStringToProducts(products)

	// add to old values new
	text := createMessageSuccessAddedProduct(newProduct)
	msg := h.createMessage(u.ChatID, text)
	if isGroup {
		list.Editors = addManyEditsProductList(u, list.Editors, len(newProduct))
		msg.ReplyMarkup = createInlineGetCurGroupList(*list.ID, *list.Name)
	} else {
		msg.ReplyMarkup = createInlineGetCurList(*list.ID, *list.Name)
	}
	list.Products = append(list.Products, newProduct...)

	// update all row of products set new products and edit many edits for editors list
	if err := h.db.ProductList().EditProductList(context.TODO(), *list); err != nil {
		return nil, err
	}

	return msg, nil
}

func (h *Hub) createMessageForEditList(ChatID int64, products, listName string, listID int, lastMsgID int, isGroup bool) *tg.EditMessageTextConfig {
	if products == emptyListMessage {
		editMsg := h.editMessage(ChatID, lastMsgID, "It seems that your list is empty 🗿, you have nothing to delete")
		editMsg.ReplyMarkup = createInlineGoToMenu()
		return editMsg
	}
	text := products + "\n\n❓❓❓\n\n" + answerEditListMessage + listName
	editMsg := h.editMessage(ChatID, lastMsgID, text)
	var userCmd TypeUserCommand
	if isGroup {
		userCmd = TypeUserCommand{
			TypeCmd: isEditGroupList,
			ListID:  &listID,
		}
	} else {
		userCmd = TypeUserCommand{
			TypeCmd: isEditList,
			ListID:  &listID,
		}
	}
	h.container.AddUserCmd(ChatID, userCmd)
	return editMsg
}

func (h *Hub) wantCompliteList(chatID int64, listName, products string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	if products == emptyListMessage {
		err := h.db.ProductList().Delete(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
		editMsg := h.editMessage(chatID, lastMsgID, isCompletesProductListMsg+listName)
		editMsg.ReplyMarkup = createInlineGoToGroups()
		if err != nil {
			return nil, err
		}
		return editMsg, nil
	}
	list, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		return nil, err
	}

	text := createMessageCompliteGroupList(*list, chatID)
	msg := h.editMessage(chatID, lastMsgID, text)

	msg.ReplyMarkup = createInlineAfterComplite(listID, *list.GroupID, listName)
	return msg, nil
}

func (h *Hub) compliteProductList(ChatID int64, name string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {

	list, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		return nil, err
	}
	if len(list.Products) == 0 {
		err := h.db.ProductList().Delete(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
		editMsg := h.editMessage(ChatID, lastMsgID, "🔴 "+*list.Name+" is deleted")
		editMsg.ReplyMarkup = createInlineGoToMenu()
		return editMsg, nil
	}
	err = h.db.ProductList().MakeListInactive(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	text := createMessageComliteUserList(*list)
	editMsg := h.editMessage(ChatID, lastMsgID, text)
	editMsg.ReplyMarkup = createInlineRecoverList(listID)
	go func() {
		if ok := h.setTimerForComliteMsg(listID); ok {
			return
		}
		err := h.db.ProductList().Delete(context.TODO(), listID)
		// What if error
		if err != nil {
			return
		}
	}()
	return editMsg, nil
}

func (h *Hub) recoverUserList(chatID int64, listID int, text string) (msg *tg.MessageConfig, err error) {
	var listName string
	// make List active in db
	if h.container.isInContainerList(listID) {
		h.container.DeleteRecoverList(listID)
		listName, err = h.db.ProductList().MakeListActive(context.TODO(), listID)

	} else {
		list := parseTextUserList(text, chatID)
		listName = *list.Name
		listID, err = h.db.ProductList().Create(context.TODO(), list)
		if err != nil {
			return nil, err
		}
	}

	msg = h.createMessage(chatID, "The list - "+listName+" is recover")
	msg.ReplyMarkup = createInlineGetCurList(listID, listName)
	return msg, nil
}

func (h *Hub) editProductList(chatID int64, listID int, indexProducts map[int]bool, isGroup bool) (*tg.MessageConfig, error) {

	list, err := h.db.ProductList().GetAllInfoProductLissIdOrName(context.TODO(), listID, "")
	if err != nil {
		return nil, err
	}

	list.Products = deleteProductByIndex(list.Products, indexProducts)

	err = h.db.ProductList().EditProductList(context.TODO(), *list)
	if err != nil {
		return nil, err
	}
	var text string
	if len(list.Products) == 0 {
		text = emptyListMessage
	} else {
		text = createMessageProductList(list.Products)
	}
	msg := h.createMessage(chatID, text)
	if isGroup {
		msg.ReplyMarkup = createInlineGetCurGroupList(*list.ID, *list.Name)
	} else {
		msg.ReplyMarkup = createInlineGetCurList(*list.ID, *list.Name)
	}
	return msg, nil
}

func (h *Hub) getNameGroupMergeList(chatID int64, listName string, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	listId, err := h.db.ProductList().GetListID(context.TODO(), listName)
	if err != nil {
		return nil, err
	}
	groups, err := h.db.ManagerGroup().UserGroup(context.TODO(), chatID)
	if err != nil {
		if err == store.NoUserGroupError {
			editMsg := h.editMessage(chatID, lastMsgID, err.Error())
			editMsg.ReplyMarkup = createInlineGoToMenu()
			return editMsg, nil
		}
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, choiceWhatGroupMerge)
	editMsg.ReplyMarkup = createInlineMergeListGroup(groups, listId)
	return editMsg, nil
}

func (h *Hub) mergeListWithGroup(chatID int64, groupID, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	err := h.db.ProductList().MergeListGroup(context.TODO(), listID, groupID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, successMergeListGroupMgs)

	editMsg.ReplyMarkup = createInlineGetCurGroup(groupID)
	return editMsg, nil
}

func (h *Hub) getMainMenu(chatID int64, lastMsgID int) *tg.EditMessageTextConfig {
	editMsg := h.editMessage(chatID, lastMsgID, cmdMenu)
	editMsg.ReplyMarkup = createInlineGetChoiceList()
	return editMsg
}
