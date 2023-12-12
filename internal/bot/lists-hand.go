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
	id, err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", *list.Name))
	msg.ReplyMarkup = createInlineGetCurList(id)
	return msg, nil
}

// Info: Select all users lists and create inline keyboard with its
func (h *Hub) getListName(chatID int64, lastMsgID int) (editMsg *tg.EditMessageTextConfig, err error) {
	lists, err := h.db.ProductList().GetAllNames(context.TODO(), chatID)
	if err != nil {
		if err == store.NoRowListOfProductError {
			editMsg = h.editMessage(chatID, lastMsgID, "Nothing is found. Create Youre First list!")
			editMsg.ReplyMarkup = createInlineNoSoloList()
			return editMsg, nil
		}
		return nil, err
	}
	keyboard := createListProductInline(lists)
	editMsg = h.editMessage(chatID, lastMsgID, listsProductsMsgHelp)

	editMsg.ReplyMarkup = &keyboard
	return editMsg, nil
}

func (h *Hub) getProductListV2(ChatID int64, lastMsgID, listID int) (editMsg *tg.EditMessageTextConfig, err error) {
	prod, err := h.db.Product().GetByListID(context.TODO(), listID, 0, 100)
	if err != nil {
		return nil, err
	}
	if len(prod) == 0 {
		editMsg = h.editMessage(ChatID, lastMsgID, emptyListMessage)
	} else {

		text := createMessageProductList(prod)
		editMsg = h.editMessage(ChatID, lastMsgID, text)
	}
	editMsg.ReplyMarkup = createProductsInline(listID)

	return editMsg, nil

}

func (h *Hub) wantAddNewProduct(chatID int64, products string, listID, lastMsgID int, isGroup bool) *tg.EditMessageTextConfig {
	var text string
	if products == emptyListMessage {
		text = addNewProductMessageReply
	} else {
		text = products + "\n❓\n" + addNewProductMessageReply
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

func (h *Hub) addNewProductV2(u store.User, products string, listID int, isGroup bool) (*tg.MessageConfig, error) {
	// FIXME: move it out
	newProduct := parseStringToProducts(products)
	var r []store.Product
	for _, p := range newProduct {
		o := store.Product{
			Product: p,
			UserID:  u.ChatID,
			ListID:  &listID,
		}
		r = append(r, o)
	}
	err := h.db.Product().Add(r)
	if err != nil {
		return nil, err
	}
	text := createMessageSuccessAddedProduct(newProduct)
	msg := h.createMessage(u.ChatID, text)
	if isGroup {
		msg.ReplyMarkup = createInlineGetCurGroupList(listID)
	} else {
		msg.ReplyMarkup = createInlineGetCurList(listID)
	}
	return msg, nil
}

func (h *Hub) compliteProductList(ChatID int64, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {
	infoList, err := h.db.ProductList().GetFoolInfoProdList(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	if len(infoList.Products) == 0 {
		err := h.db.ProductList().Delete(context.TODO(), listID)
		if err != nil {
			return nil, err
		}
		editMsg := h.editMessage(ChatID, lastMsgID, "🔴 List is deleted")
		editMsg.ReplyMarkup = createInlineGoToMenu()
		return editMsg, nil
	}
	err = h.db.ProductList().MakeListInactive(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	text := createMessageComliteUserList(infoList.Products, *infoList.List.Name)
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
		splitedText := splitText(text, '\n')

		listName = parseNameTextList(splitedText)

		list := &store.ProductList{
			Name:    &listName,
			OwnerID: &chatID,
		}
		listID, err = h.db.ProductList().Create(context.TODO(), list)
		if err != nil {
			return nil, err
		}
		prod := parseTextToProd(splitedText, chatID, listID)
		if err = h.db.Product().Add(prod); err != nil {
			return nil, err
		}
	}

	msg = h.createMessage(chatID, "The list - "+listName+" is recover 👀")
	msg.ReplyMarkup = createInlineGetCurList(listID)
	return msg, nil
}

func (h *Hub) wantDeleteProd(chatID int64, listID, offset, lastMsgID int, isGroup bool) (*tg.EditMessageTextConfig, error) {
	count, err := h.db.Product().CountByListID(context.TODO(), listID)
	if err != nil {
		// TODO: Make message if no exist products
		if err == store.NoProuductExistError {
			return nil, err
		}
		return nil, err
	}
	prod, err := h.db.Product().GetByListID(context.TODO(), listID, offset, 10)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "Select product to delete it :-)")

	editMsg.ReplyMarkup = createInlineEditProdKb(prod, listID, offset, 10, count, isGroup)

	return editMsg, nil
}

func (h *Hub) deleteProd(chatID int64, prodID, lastMsgID int, sListID string, isGroup bool) (*tg.EditMessageTextConfig, error) {
	err := h.db.Product().Delete(context.TODO(), prodID)
	if err != nil {
		return nil, err
	}
	editMsg := h.editMessage(chatID, lastMsgID, "🆑 Product success deleted")
	editMsg.ReplyMarkup = createInlineGetProdDel(sListID, isGroup)
	return editMsg, nil
}

func (h *Hub) getNameGroupMergeList(chatID int64, listID, lastMsgID int) (*tg.EditMessageTextConfig, error) {

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
	editMsg.ReplyMarkup = createInlineMergeListGroup(groups, listID)
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
