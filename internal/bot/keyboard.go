package telegram

import (
	"fmt"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// keyboard for generatins
var menuKeyboard tgbotapi.ReplyKeyboardMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(buttonListList),
		tgbotapi.NewKeyboardButton(buttonCreateList),
		tgbotapi.NewKeyboardButton(buttonNewGroup),
	),
)

func createInlineMarkup(lists []store.ProductList) tgbotapi.InlineKeyboardMarkup {
	var listOfProductList tgbotapi.InlineKeyboardMarkup
	var button tgbotapi.InlineKeyboardButton
	var buttonRow []tgbotapi.InlineKeyboardButton
	for _, list := range lists {
		CallBack := fmt.Sprintf("%d", list.ID)
		button = tgbotapi.InlineKeyboardButton{Text: fmt.Sprintf("%d", list.ID), CallbackData: &CallBack}
		buttonRow = append(buttonRow, button)
		listOfProductList.InlineKeyboard = append(listOfProductList.InlineKeyboard, buttonRow)
	}

	return listOfProductList
}
