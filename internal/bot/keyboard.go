package telegram

import (
	"strconv"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// keyboard for generatins
var (
	menuKeyboard = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton(buttonListList),
			tg.NewKeyboardButton(buttonCreateList),
			tg.NewKeyboardButton(buttonNewGroup),
			tg.NewKeyboardButton(buttonGetUserGroup),
		),
	)
)

// Info: creater keyboard on bottom for list products

func createListProductInline(lists []store.ProductList) tg.InlineKeyboardMarkup {
	var listOfProductList tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		callBack := createCallBackListProducts(*list.ID, *list.Name)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		buttonRow := []tg.InlineKeyboardButton{button}
		listOfProductList.InlineKeyboard = append(listOfProductList.InlineKeyboard, buttonRow)
	}
	return listOfProductList
}

func createProductsInline(listName string) *tg.InlineKeyboardMarkup {
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("add", "add-"+listName),
			tg.NewInlineKeyboardButtonData("change", "change-"+listName),
			tg.NewInlineKeyboardButtonData("complite", "comple-"+listName),
		),
	)
	return &keyboard
}

func createInlineGetCurList(listID int, listName string) tg.InlineKeyboardMarkup {
	data := createCallBackListProducts(listID, listName)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("View list", *data),
		),
	)
	return keyboard
}

func createInlineGroupName(groups []store.GroupInfo) *tg.InlineKeyboardMarkup {
	var keyboard = tg.InlineKeyboardMarkup{}
	var groupButton tg.InlineKeyboardButton
	for _, group := range groups {
		callBack := createCallBackGroupLists(group.ID)
		groupButton = tg.InlineKeyboardButton{Text: group.GroupName, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{groupButton},
		)
	}
	return &keyboard
}

func createInlineGroupList(lists []store.ProductList, groupID int, isOwnerGroup bool) tg.InlineKeyboardMarkup {
	keyboard := tg.InlineKeyboardMarkup{}
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		callBack := createCallBackListProducts(*list.ID, *list.Name)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}
	if isOwnerGroup {
		row := createInlineGroupActions(groupID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return keyboard
}

func createInlineGroupActions(groupID int) []tg.InlineKeyboardButton {
	sGroupID := strconv.Itoa(groupID)
	buttonRow := tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("create", "createGroupList-"+sGroupID),
		tg.NewInlineKeyboardButtonData("add users", "addUserGroup-"+sGroupID),
		tg.NewInlineKeyboardButtonData("del users", "GetDeleteUser-"+sGroupID),
	)
	return buttonRow
}

func createInlineDeleteUser(users []store.User, groupID int) *tg.InlineKeyboardMarkup {
	keyboard := &tg.InlineKeyboardMarkup{}
	var button tg.InlineKeyboardButton
	for _, user := range users {
		sUsID, sGrID := strconv.Itoa(user.ChatID), strconv.Itoa(groupID)
		callBack := createCallBackDeleteUserGroup(sUsID, sGrID)
		button = tg.NewInlineKeyboardButtonData(user.UserName, *callBack)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{button})
	}
	return keyboard
}
