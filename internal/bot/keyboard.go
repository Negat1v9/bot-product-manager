package telegram

import (
	"strconv"

	"github.com/Negat1v9/telegram-bot-orders/store"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Info: inline keyboard for user choice list or group-list
func createInlineGetChoiceList() *tg.InlineKeyboardMarkup {
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
	return &kb
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

func createProductsInline(listName string, listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⭐ Add", *createCallBackOneParam(prefixAddProductList, listName)),
			tg.NewInlineKeyboardButtonData(
				"🆑 Delete", *createCallBackOneParam(prefixChangeList, listName)),
			tg.NewInlineKeyboardButtonData(
				"✅ Complite", *createCallBackFewParam(prefixCompliteSoloList, sListID, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔓 Attach to another group", *createCallBackOneParam(prefixToMergeListGroup, listName),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🪅 Menu", prefixGetMainMenu),
		),
	)
	return &keyboard
}

func createInlineProductsGroup(listName string, listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⭐ Add", *createCallBackOneParam(prefixAddProductGroup, listName)),
			tg.NewInlineKeyboardButtonData(
				"🆑 Delete", *createCallBackOneParam(prefixChangeGroupList, listName)),
			tg.NewInlineKeyboardButtonData(
				"✅ Complite", *createCallBackFewParam(prefixWantCompliteList, sListID, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🪅 Menu", prefixGetMainMenu),
		),
	)
	return &keyboard
}

func createInlineAfterListCreated(groupID int, newListID int) *tg.InlineKeyboardMarkup {
	sGroupID, sNewListID := strconv.Itoa(groupID), strconv.Itoa(newListID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"📰 Connect", *createCallBackFewParam(prefixWantConnectTemplate, sGroupID, sNewListID)),
			tg.NewInlineKeyboardButtonData(
				"👀 Refuse", *createCallBackOneParam(prefixCallBackListGroup, sGroupID)),
		),
	)
	return &kb
}

func createInlineAfterComplite(listID int, listName string) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"📩 Save", *createCallBackOneParam(prefixSaveAsTemplete, sListID)),
			tg.NewInlineKeyboardButtonData(
				"🎯 Complite", *createCallBackFewParam(prefixCompliteList, sListID, listName)),
		),
	)
	return &kb
}

func createInlineGetCurList(listID int, listName string) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	data := createCallBackFewParam(prefixCallBackListProduct, sListID, listName)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📕 View list", *data),
		),
	)
	return &keyboard
}

func createInlineGetCurGroup(groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	data := createCallBackOneParam(prefixCallBackListGroup, sGroupID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔍 See the group", *data,
			),
		),
	)
	return &kb
}

func creaetInlineBackToGroupButton(groupID int) *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"👈 Return to group", *createCallBackOneParam(
					prefixCallBackListGroup, strconv.Itoa(groupID)),
			),
		),
	)
	return &kb
}

func createInlineGoToGroups() *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔍 View groups", prefixGetGroupLists)),
	)
	return &kb
}

func createInlineGoToMenu() *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⚡ Menu", prefixGetMainMenu)),
	)
	return &kb
}

func createInlineMakeSureDelete(groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"Refuse 🛑", prefixGetGroupLists),
			tg.NewInlineKeyboardButtonData(
				"Confirm ✅", *createCallBackOneParam(prefixLeaveOwnerGroup, sGroupID)),
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
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"💴 Menu", prefixGetMainMenu),
		))
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
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🌵 Back", prefixGetUserList)))
	return &keyboard
}

func createInlineGroupList(lists []store.ProductList, groupID int) *tg.InlineKeyboardMarkup {
	var keyboard tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		sListID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixCallBackGroupProductList, sListID, *list.Name)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}

	row := createInlineGroupActions(groupID)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row...)
	return &keyboard
}

func createInlineGroupActions(groupID int) [][]tg.InlineKeyboardButton {
	sGroupID := strconv.Itoa(groupID)
	rows := [][]tg.InlineKeyboardButton{
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				choiceCreateSoloList, *createCallBackOneParam(prefixCreateGroupList, sGroupID)),
			tg.NewInlineKeyboardButtonData(
				choiceGetAllUsersGroup, *createCallBackOneParam(prefixGetAllUsersGroup, sGroupID)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"📬 Templates", *createCallBackOneParam(prefixGetGroupTemplates, sGroupID),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🏇 Return to groups", prefixGetGroupLists,
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"💔 Leave group", *createCallBackOneParam(prefixLeaveGroup, sGroupID)),
		),
	}

	return rows
}

func createInlineGroupListTemplates(lists []store.ProductList, groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	var keyboard tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		sListID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixGetOneTemplate, sListID, *list.Name, sGroupID)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🦀 Returns to group", *createCallBackOneParam(prefixCallBackListGroup, sGroupID))))

	return &keyboard
}

func createInlineTemplateForConnect(lists []store.ProductList, groupID, newListID int) *tg.InlineKeyboardMarkup {
	sGroupID, sNewListID := strconv.Itoa(groupID), strconv.Itoa(newListID)
	var kb tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		sListID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixGetListForTemplateMerge, sListID, sNewListID, sGroupID)
		button = tg.InlineKeyboardButton{Text: *list.Name, CallbackData: callBack}
		kb.InlineKeyboard = append(kb.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}
	kb.InlineKeyboard = append(kb.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🐣 Cancel", *createCallBackOneParam(prefixCallBackListGroup, sGroupID))))
	return &kb
}

func createInlineConnectTemplate(listID int, sNewListID, sGroupID string) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"➕ Connect", *createCallBackFewParam(prefixConnectTemplate, sListID, sNewListID, sGroupID))),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"👈 Back", *createCallBackFewParam(prefixWantConnectTemplate, sGroupID, sNewListID))),
	)
	return &kb
}

func createInlineTemplateActions(listID int, listName, sGroupID string) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⭐ Add product", *createCallBackOneParam(prefixAddProductList, listName)),
			tg.NewInlineKeyboardButtonData(
				"🚮 Remove product", *createCallBackOneParam(prefixChangeList, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🆑 Delete template", *createCallBackFewParam(prefixCompliteList, sListID, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"👈 Back to group", *createCallBackOneParam(prefixCallBackListGroup, sGroupID)),
		),
	)
	return &kb
}

func creaetInlineUsersGroupActions(groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🟢 Invite", *createCallBackOneParam(prefixAddUserGroup, sGroupID)),
			tg.NewInlineKeyboardButtonData(
				"🚫 Delete", *createCallBackOneParam(prefixGetUserToDelete, sGroupID)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"👈 Return to group", *createCallBackOneParam(prefixCallBackListGroup, sGroupID)),
		),
	)
	return &kb
}

func createInlineDeleteUser(users []store.User, groupID int, ownerId int64) *tg.InlineKeyboardMarkup {
	keyboard := &tg.InlineKeyboardMarkup{}
	var button tg.InlineKeyboardButton
	sGroupID := strconv.Itoa(groupID)
	for _, user := range users {
		// skip button for delete ownerGroup
		if user.ChatID == ownerId {
			continue
		}
		sUsID := strconv.FormatInt(user.ChatID, 10)
		callBack := createCallBackFewParam(prefixCallBackDelUserFromGr, sUsID, sGroupID)
		button = tg.NewInlineKeyboardButtonData(createButtonDeleteUser(*user.UserName), *callBack)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{button})
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("👈 Return to group", *createCallBackOneParam(prefixCallBackListGroup, sGroupID))))
	return keyboard
}

func createInlineInviteUserGroup(groupID int, newUserID int64) *tg.InlineKeyboardMarkup {
	sUsID, sGrID := strconv.FormatInt(newUserID, 10), strconv.Itoa(groupID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"💚 Join",
				*createCallBackFewParam(prefixCallBackInsertUserGroup, sUsID, sGrID)),
			tg.NewInlineKeyboardButtonData(
				"💔 Refuse",
				*createCallBackFewParam(prefixCallBackRefuseUserGroup, sUsID, sGrID)),
		),
	)
	return &keyboard
}
