package telegram

import (
	"context"
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: create message with information to user, about he have to forwad its
func (h *Hub) answerToCreateList(ChatID int64) *tgbotapi.MessageConfig {
	msg := h.createMessage(ChatID, answerCreateListMsg)
	return msg
}

// Info: Create list in database, and answer with new list name
func (h *Hub) createList(ChatID int64, list *store.ProductList) (*tgbotapi.MessageConfig, error) {

	err := h.db.ProductList().Create(context.TODO(), list)
	if err != nil {
		return nil, err
	}
	msg := h.createMessage(ChatID, fmt.Sprintf("New list %s is created success", list.Name))
	return msg, nil
}

// Info: Select all users lists and create inline keyboard with its
func (h *Hub) getListName(UserID, ChatID int64) (msg *tgbotapi.MessageConfig, err error) {
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

func (h *Hub) getProductList(ChatID int64, listID int, listName string) (msg *tgbotapi.MessageConfig, err error) {
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

func (h *Hub) addNewProduct(ChatID int64, p store.Product, listName string) (*tgbotapi.MessageConfig, error) {
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
