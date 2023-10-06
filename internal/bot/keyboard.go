package telegram

import (
	"github.com/Negat1v9/telegram-bot-orders/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// keyboard for generatins
var (
	menuKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonListList),
			tgbotapi.NewKeyboardButton(buttonCreateList),
			tgbotapi.NewKeyboardButton(buttonNewGroup),
		),
	)
)

// Info: creater keyboard on bottom for list products
func createProductsInline(listName string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выполнить", "complete-"+listName),
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "add-"+listName),
			tgbotapi.NewInlineKeyboardButtonData("Изменить", "change-"+listName),
		),
	)
	return &keyboard
}

func createListProductInline(lists []store.ProductList) tgbotapi.InlineKeyboardMarkup {
	var listOfProductList tgbotapi.InlineKeyboardMarkup
	var button tgbotapi.InlineKeyboardButton
	for _, list := range lists {
		CallBack := createCallBackListProducts(list.ID, list.Name)
		button = tgbotapi.InlineKeyboardButton{Text: list.Name, CallbackData: CallBack}
		buttonRow := []tgbotapi.InlineKeyboardButton{button}
		listOfProductList.InlineKeyboard = append(listOfProductList.InlineKeyboard, buttonRow)
	}
	return listOfProductList
}

func createInlineGetCurList(listID int, listName string) tgbotapi.InlineKeyboardMarkup {
	data := createCallBackListProducts(listID, listName)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Смотреть лист", *data),
		),
	)
	return keyboard
}
