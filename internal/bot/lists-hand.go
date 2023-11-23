package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Hub) createMsgToCreateList(ChatID int64, lastMsgID int) *tg.EditMessageTextConfig {
	editMsg := h.editMessage(ChatID, lastMsgID, answerCreateListMsg)
	return editMsg
}

// Info: Create list in database, and answer with new list name
func (h *Hub) createList(ChatID int64, list *store.ProductList) (*tg.MessageConfig, error) {
	clearName := makeNameClear(*list.Name)
	list.Name = &clearName
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
	lists, err := h.db.ProductList().GetAll(context.TODO(), chatID)
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
	product, err := h.db.Product().GetAllProducts(context.TODO(), listID)
	if err != nil {
		if err == store.NoRowProductError {
			editMsg = h.editMessage(ChatID, lastMsgID, emptyListMessage)

		} else {
			return nil, err
		}

	} else if len(product.Products) == 0 {
		editMsg = h.editMessage(ChatID, lastMsgID, emptyListMessage)
	} else {
		text := createMessageProductList(product.Products)
		editMsg = h.editMessage(ChatID, lastMsgID, text)
	}

	editMsg.ReplyMarkup = createProductsInline(listName, listID)
	return editMsg, nil
}

func (h *Hub) wantAddNewProduct(chatID int64, products, listName string, lastMsgID int, isGroup bool) *tg.EditMessageTextConfig {
	var text string
	if products == emptyListMessage {
		text = addNewProductMessageReply + listName
	} else {
		text = products + "\n❓\n" + addNewProductMessageReply + listName
	}
	var editMsg *tg.EditMessageTextConfig

	if isGroup {
		text = products + "\n❓\n" + addNewProductAtGroupList + listName
		editMsg = h.editMessage(chatID, lastMsgID, text)
	} else {
		editMsg = h.editMessage(chatID, lastMsgID, text)
	}
	return editMsg
}

// TODO: Make something with length func
func (h *Hub) addNewProduct(u store.User, Products, listName string, isGroup bool) (*tg.MessageConfig, error) {
	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
	if err != nil {
		return nil, err
	}
	newProduct := parseStringToProducts(Products, listID)

	lastProduct, err := h.db.Product().GetAllProducts(context.TODO(), listID)
	if err != nil {
		// if not row exist
		if err == store.NoRowProductError {
			err = h.db.Product().Create(context.TODO(), &store.Product{ListID: listID})
			lastProduct = &store.Product{
				ListID:   listID,
				Products: []string{},
				Editors:  []store.Editors{},
			}
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	text := createMessageSuccessAddedProduct(newProduct.Products)
	msg := h.createMessage(u.ChatID, text)
	if isGroup {
		newProduct.Editors = addManyEditsProductList(u, lastProduct.Editors, len(newProduct.Products))
		msg.ReplyMarkup = createInlineGetCurGroupList(listID, listName)
	} else {
		msg.ReplyMarkup = createInlineGetCurList(listID, listName)
	}
	// add to old values new
	if len(lastProduct.Products) > 0 {
		lastProduct.Products = append(lastProduct.Products, newProduct.Products...)
		newProduct.Products = lastProduct.Products
	}
	// update all row of products set new products and edit many edits for editors list
	if err := h.db.Product().Add(context.TODO(), newProduct); err != nil {
		return nil, err
	}

	return msg, nil
}

func (h *Hub) createMessageForEditList(ChatID int64, products, listName string, lastMsgID int, isGroup bool) *tg.EditMessageTextConfig {
	if products == emptyListMessage {
		editMsg := h.editMessage(ChatID, lastMsgID, "It seems that your list is empty 🗿, you have nothing to delete")
		editMsg.ReplyMarkup = createInlineGoToMenu()
		return editMsg
	}
	var editMsg *tg.EditMessageTextConfig
	if isGroup {
		text := products + "\n\n❓❓❓\n\n" + answerEditGroupList + listName
		editMsg = h.editMessage(ChatID, lastMsgID, text)
	} else {
		text := products + "\n\n❓❓❓\n\n" + answerEditListMessage + listName
		editMsg = h.editMessage(ChatID, lastMsgID, text)
	}
	return editMsg
}

func (h *Hub) wantCompliteList(chatID int64, listName, products string, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	if products == emptyListMessage {
		err := h.db.ProductList().Delete(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
		editMsg := h.editMessage(chatID, lastMsgID, isCompletesProductListMsg)
		if err != nil {
			return nil, err
		}
		return editMsg, nil
	}
	productsInfo, err := h.db.Product().GetAllProducts(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	// get groupID from db
	productList, err := h.db.ProductList().GetGroupListByID(context.TODO(), productsInfo.ListID)
	if err != nil {
		return nil, err
	}
	// editorString :=
	text := createMessageCompliteList(*productsInfo, chatID)
	msg := h.editMessage(chatID, lastMsgID, text)

	msg.ReplyMarkup = createInlineAfterComplite(listID, *productList.GroupID, listName)
	return msg, nil
}

func (h *Hub) compliteProductList(ChatID int64, name, sGrID, text string, listID, lastMsgID int) (*tg.EditMessageReplyMarkupConfig, error) {

	// err := h.db.ProductList().Delete(context.TODO(), listID)
	err := h.db.ProductList().MakeListInactive(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	var markup *tg.InlineKeyboardMarkup

	if sGrID == "" {
		markup = createInlineRecoverList(listID)
	} else {
		groupID := convSToI[int](sGrID, 0)
		go h.sendComplitedListGroupDelay(listID, groupID, text)
		markup = createInlineRecoverGroupList(listID, sGrID, name)
	}

	msg := h.editReplyMarkup(ChatID, markup, lastMsgID)
	return msg, nil
}

func (h *Hub) makeListActive(chatID int64, listID int) (*tg.MessageConfig, error) {
	listName, err := h.db.ProductList().MakeListActive(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(chatID, "The list - "+listName+" is recover")
	msg.ReplyMarkup = createInlineGetCurList(listID, listName)
	return msg, nil
}

func (h *Hub) editProductList(chatID int64, listName string, indexProducts map[int]bool, isGroup bool) (*tg.MessageConfig, error) {
	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
	if err != nil {
		return nil, err
	}
	products, err := h.db.Product().GetAllProducts(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	products.Products = deleteProductByIndex(products.Products, indexProducts)

	err = h.db.Product().Add(context.TODO(), *products)
	if err != nil {
		return nil, err
	}
	var text string
	if len(products.Products) == 0 {
		text = emptyListMessage
	} else {
		text = createMessageProductList(products.Products)
	}
	msg := h.createMessage(chatID, text)
	if isGroup {
		msg.ReplyMarkup = createInlineGetCurGroupList(listID, listName)
	} else {
		msg.ReplyMarkup = createInlineGetCurList(listID, listName)
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
