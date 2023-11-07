package telegram

import (
	"strconv"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: inline keyboard for user choice list or group-list
func createInlineGetChoiceList() tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				choiceUserList, prefixGetUserList),
			tg.NewInlineKeyboardButtonData(
				choiceGroupList, prefixGetGroupLists),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				choiceCreateSoloList, prefixCreateSoloList),
			tg.NewInlineKeyboardButtonData(
				choiceCreateGroup, prefixCreateGroup),
		),
	)
	return kb
}

// Info: creater keyboard on bottom for list products
func createListProductInline(lists []store.ProductList) tg.InlineKeyboardMarkup {
	var listOfProductList tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		stID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixCallBackListProduct, stID, *list.Name)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		buttonRow := []tg.InlineKeyboardButton{button}
		listOfProductList.InlineKeyboard = append(listOfProductList.InlineKeyboard, buttonRow)
	}
	return listOfProductList
}

func createProductsInline(listName string) *tg.InlineKeyboardMarkup {
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"add", *createCallBackOneParam(prefixAddProductList, listName)),
			tg.NewInlineKeyboardButtonData(
				"change", *createCallBackOneParam(prefixChangeList, listName)),
			tg.NewInlineKeyboardButtonData(
				"complite", *createCallBackOneParam(prefixCompliteList, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔓 merge with group", *createCallBackOneParam(prefixToMergeListGroup, listName),
			),
		),
	)
	return &keyboard
}

func createInlineGetCurList(listID int, listName string) tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	data := createCallBackFewParam(prefixCallBackListProduct, sListID, listName)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("View list", *data),
		),
	)
	return keyboard
}

func createInlineGetCurGroup(groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	data := createCallBackOneParam(prefixCallBackListGroup, sGroupID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"See the group? 🔍", *data,
			),
		),
	)
	return &kb
}

func createInlineGroupName(groups []store.GroupInfo) *tg.InlineKeyboardMarkup {
	var keyboard = tg.InlineKeyboardMarkup{}
	var groupButton tg.InlineKeyboardButton
	for _, group := range groups {
		sGroupID := strconv.Itoa(group.ID)
		callBack := createCallBackOneParam(prefixCallBackListGroup, sGroupID)
		groupButton = tg.InlineKeyboardButton{Text: group.GroupName, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{groupButton},
		)
	}
	return &keyboard
}

// prefix-GroupID-ListID
func createInlineMergeListGroup(groups []store.GroupInfo, listID int) *tg.InlineKeyboardMarkup {
	var keyboard = tg.InlineKeyboardMarkup{}
	var groupButton tg.InlineKeyboardButton
	sListID := strconv.Itoa(listID)
	for _, group := range groups {
		sGroupID := strconv.Itoa(group.ID)
		callBack := createCallBackFewParam(prefixMergeListWithGroup, sGroupID, sListID)
		groupButton = tg.InlineKeyboardButton{Text: group.GroupName, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{groupButton},
		)
	}
	return &keyboard
}

func createInlineGroupList(lists []store.ProductList, groupID int, isOwnerGroup bool) *tg.InlineKeyboardMarkup {
	var keyboard tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		sListID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixCallBackListProduct, sListID, *list.Name)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}
	if isOwnerGroup {
		row := createInlineGroupActions(groupID)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return &keyboard
}

func createInlineGroupActions(groupID int) []tg.InlineKeyboardButton {
	sGroupID := strconv.Itoa(groupID)
	buttonRow := tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(
			choiceCreateSoloList, *createCallBackOneParam(prefixCreateGroupList, sGroupID)),
		tg.NewInlineKeyboardButtonData(
			choiceGetAllUsersGroup, *createCallBackOneParam(prefixGetAllUsersGroup, sGroupID)),
	)
	return buttonRow
}

func creaetInlineUsersGroupActions(groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	kb := tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(
			"🟢 invite", *createCallBackOneParam(prefixAddUserGroup, sGroupID)),
		tg.NewInlineKeyboardButtonData(
			"🚫 delete", *createCallBackOneParam(prefixGetUserToDelete, sGroupID)),
	),
	)
	return &kb
}

func createInlineDeleteUser(users []store.User, groupID int, ownerId int64) *tg.InlineKeyboardMarkup {
	keyboard := &tg.InlineKeyboardMarkup{}
	var button tg.InlineKeyboardButton
	for _, user := range users {
		// skip button for delete ownerGroup
		if user.ChatID == ownerId {
			continue
		}
		sUsID, sGrID := strconv.FormatInt(user.ChatID, 10), strconv.Itoa(groupID)
		callBack := createCallBackFewParam(prefixCallBackDelUserFromGr, sUsID, sGrID)
		button = tg.NewInlineKeyboardButtonData(createButtonDeleteUser(*user.UserName), *callBack)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{button})
	}
	return keyboard
}

func createInlineInviteUserGroup(groupID int, newUserID int64) *tg.InlineKeyboardMarkup {
	sUsID, sGrID := strconv.FormatInt(newUserID, 10), strconv.Itoa(groupID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"join",
				*createCallBackFewParam(prefixCallBackInsertUserGroup, sUsID, sGrID)),
			tg.NewInlineKeyboardButtonData(
				"refuse",
				*createCallBackFewParam(prefixCallBackRefuseUserGroup, sUsID, sGrID)),
		),
	)
	return &keyboard
}
