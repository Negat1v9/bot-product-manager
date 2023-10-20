package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
func (h *Hub) getListName(chatID int64, lastMsgID int) (editMsg *tg.EditMessageTextConfig, err error) {
	lists, err := h.db.ProductList().GetAll(context.TODO(), int(chatID))
	if err != nil {
		if err == store.NoRowListOfProductError {
			editMsg = h.editMessage(chatID, lastMsgID, "Nothing is found. Create Youre First list!")
			// editMsg = h.createMessage(chatID, "Nothing is found. Create Youre First list!")
			return editMsg, nil
		}
		return nil, err
	}
	keyboard := createListProductInline(lists)
	editMsg = h.editMessage(chatID, lastMsgID, "List of Product-lists")
	// editMsg = h.createMessage(chatID, "List of Product-lists")
	editMsg.ReplyMarkup = &keyboard
	return editMsg, nil
}

func (h *Hub) getChoiceLists(chatID int64) *tg.MessageConfig {
	msg := h.createMessage(chatID, "select list")
	msg.ReplyMarkup = createInlineGetChoiceList()
	return msg
}

func (h *Hub) getProductList(ChatID int64, lastMsgID, listID int, listName string) (*tg.EditMessageTextConfig, error) {
	product, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		if err == store.NoRowProductError {
			editMsg := h.editMessage(ChatID, lastMsgID, emptyListMessage)
			editMsg.ReplyMarkup = createProductsInline(listName)
			return editMsg, nil

		}
		return nil, err
	}
	text := createMessageProductList(product.Products)
	editMsg := h.editMessage(ChatID, lastMsgID, text)

	editMsg.ReplyMarkup = createProductsInline(listName)
	return editMsg, nil
}

func (h *Hub) addNewProduct(ChatID int64, Products, listName string) (*tg.MessageConfig, error) {
	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
	if err != nil {
		return nil, err
	}
	newProduct := parseStringToProducts(Products, listID)

	lastProduct, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		// if not row exist
		if err == store.NoRowProductError {
			err = h.db.Product().Create(context.TODO(), listID)
			lastProduct = &store.Product{
				ListID:   listID,
				Products: []string{},
			}
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	text := createMessageSuccessAddedProduct(newProduct.Products)
	// add to old values new
	if len(lastProduct.Products) > 0 {
		lastProduct.Products = append(lastProduct.Products, newProduct.Products...)
		newProduct.Products = lastProduct.Products
	}
	if err := h.db.Product().Add(context.TODO(), newProduct); err != nil {
		return nil, err
	}

	msg := h.createMessage(ChatID, text)
	msg.ReplyMarkup = createInlineGetCurList(listID, listName)
	return msg, nil
}

func (h *Hub) createMessageForEditList(ChatID int64, listName string) *tg.MessageConfig {
	msg := h.createMessage(ChatID, answerEditListMessage+listName)
	return msg
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

func (h *Hub) editProductList(chatID int64, listName string, indexProducts map[int]bool) (*tg.MessageConfig, error) {
	listID, err := h.db.ProductList().GetListID(context.TODO(), listName)
	if err != nil {
		return nil, err
	}
	products, err := h.db.Product().GetAll(context.TODO(), listID)
	if err != nil {
		return nil, err
	}
	products.Products = deleteProductByIndex(products.Products, indexProducts)

	err = h.db.Product().Add(context.TODO(), *products)
	if err != nil {
		// TODO: Not deleted
		return nil, err
	}
	text := createMessageProductList(products.Products)
	msg := h.createMessage(chatID, text)
	msg.ReplyMarkup = createProductsInline(listName)
	return msg, nil
}
