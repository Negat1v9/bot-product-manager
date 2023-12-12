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
	)
	return &kb
}

// Info: creater keyboard on bottom for list products
func createListProductInline(lists []store.ProductList) tg.InlineKeyboardMarkup {
	var kb tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		stID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixCallBackListProduct, stID)
		button = tg.InlineKeyboardButton{Text: "📜 " + *list.Name, CallbackData: callBack}
		buttonRow := []tg.InlineKeyboardButton{button}
		kb.InlineKeyboard = append(kb.InlineKeyboard, buttonRow)
	}
	kb.InlineKeyboard = append(kb.InlineKeyboard, tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData(
			"New-list 📚", prefixCreateSoloList),
	))
	return kb
}

func createProductsInline(listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⭐ Add", *createCallBackOneParam(prefixAddProductList, sListID)),
			tg.NewInlineKeyboardButtonData(
				"🆑 Delete", *createCallBackFewParam(prefixGetPageProdDelete, sListID, "0")),
			tg.NewInlineKeyboardButtonData(
				"✅ Complite", *createCallBackOneParam(prefixCompliteSoloList, sListID)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔓 Attach to another group", *createCallBackOneParam(prefixToMergeListGroup, sListID),
			),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🪅 Menu", prefixGetMainMenu),
		),
	)
	return &keyboard
}

func createInlineProductsGroup(listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⭐ Add", *createCallBackOneParam(prefixAddProductGroup, sListID)),
			tg.NewInlineKeyboardButtonData(
				"🆑 Delete", *createCallBackFewParam(prefixGetPageGroupProdDelete, sListID, "0")),
			tg.NewInlineKeyboardButtonData(
				"✅ Complite", *createCallBackOneParam(prefixWantCompliteGrList, sListID)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🪅 groups", prefixGetGroupLists),
		),
	)
	return &keyboard
}

func createInlineAfterComplite(listID, groupID int, listName string) *tg.InlineKeyboardMarkup {
	sListID, sGroupID := strconv.Itoa(listID), strconv.Itoa(groupID)
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🎯 Complite", *createCallBackFewParam(prefixCompliteList, sListID, sGroupID, listName)),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⏪ Back", *createCallBackOneParam(prefixCallBackGroupProductList, sListID)),
		),
	)
	return &kb
}

func createInlineRecoverList(listID int) *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔮 Recover", *createCallBackFewParam(prefixRestoreList, strconv.Itoa(listID))),
		),
	)
	return &kb
}

func createInlineRecoverGroupList(listID int, sGroupID, listName string) *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🔮 Recover", *createCallBackFewParam(prefixRestoreGroupList, strconv.Itoa(listID), sGroupID, listName)),
		),
	)
	return &kb
}

func createInlineGetCurList(listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	data := createCallBackFewParam(prefixCallBackListProduct, sListID)
	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📕 View list", *data),
		),
	)
	return &keyboard
}

func createInlineGetCurGroupList(listID int) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	data := createCallBackFewParam(prefixCallBackGroupProductList, sListID)
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

func createInlineNoSoloList() *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"📚 New-list", prefixCreateSoloList)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⚡ Menu", prefixGetMainMenu)),
	)
	return &kb
}
func createInlineNoGroupUser() *tg.InlineKeyboardMarkup {
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🥷 New-group", prefixCreateGroup)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⚡ Menu", prefixGetMainMenu)),
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
		groupButton = tg.InlineKeyboardButton{Text: "👥 " + group.GroupName, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			[]tg.InlineKeyboardButton{groupButton},
		)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"💴 Menu", prefixGetMainMenu)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🥷 New-group", prefixCreateGroup)),
	)
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
		groupButton = tg.InlineKeyboardButton{Text: "👥 " + group.GroupName, CallbackData: callBack}
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

func createInlineEditProdKb(prod []store.Product, listID, offset, limit, manyOfItem int, isGroup bool) *tg.InlineKeyboardMarkup {
	sListID := strconv.Itoa(listID)
	var pref, prefBack string
	if isGroup {
		pref = prefixDeleteGrProd
		prefBack = prefixCallBackGroupProductList
	} else {
		pref = prefixDeleteProd
		prefBack = prefixCallBackListProduct
	}
	kb := tg.InlineKeyboardMarkup{}
	var button tg.InlineKeyboardButton
	for _, p := range prod {
		sProdID := strconv.Itoa(p.ID)
		cb := createCallBackFewParam(pref, sProdID, sListID)
		button = tg.InlineKeyboardButton{Text: "🔴 " + p.Product, CallbackData: cb}
		kb.InlineKeyboard = append(kb.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}
	addPaginationInPlace(&kb, offset, limit, manyOfItem, sListID)
	addButtonBack(&kb, prefBack, sListID)
	return &kb
}

func addPaginationInPlace(kb *tg.InlineKeyboardMarkup, offset, limit, manyOfItem int, sListID string) {
	var button tg.InlineKeyboardButton
	var row []tg.InlineKeyboardButton
	button = tg.NewInlineKeyboardButtonData("• 1 •",
		*createCallBackFewParam(prefixGetPageProdDelete, sListID, "0"))
	if manyOfItem <= limit {
		row = append(row, button)
		return
	}
	sPage := ""
	if manyOfItem <= 50 {
		for i := 0; i <= manyOfItem; i += limit {
			sPage = strconv.Itoa(i/limit + 1)
			sOffset := strconv.Itoa(i)
			if offset == i {
				button = tg.NewInlineKeyboardButtonData("• "+sPage+" •",
					*createCallBackFewParam(prefixGetPageProdDelete, sListID, sOffset))
			} else {
				button = tg.NewInlineKeyboardButtonData(sPage,
					*createCallBackFewParam(prefixGetPageProdDelete, sListID, sOffset))
			}
			row = append(row, button)
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}
}

func addButtonBack(kb *tg.InlineKeyboardMarkup, pref, sListID string) {
	kb.InlineKeyboard = append(kb.InlineKeyboard,
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"⏪ Back", *createCallBackOneParam(pref, sListID))))
}

func createInlineGetProdDel(sListID string, isGroup bool) *tg.InlineKeyboardMarkup {
	var prefGetToDel, prefGetList string
	if isGroup {
		prefGetToDel = prefixGetPageGroupProdDelete
		prefGetList = prefixCallBackGroupProductList
	} else {
		prefGetToDel = prefixGetPageProdDelete
		prefGetList = prefixCallBackListProduct
	}
	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"🃏 Delete more", *createCallBackFewParam(prefGetToDel, sListID, "0"))),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"✅ View List", *createCallBackOneParam(prefGetList, sListID))),
	)
	return &kb
}

func createInlineGroupList(lists []store.ProductList, groupID int) *tg.InlineKeyboardMarkup {
	sGroupID := strconv.Itoa(groupID)
	var keyboard tg.InlineKeyboardMarkup
	var button tg.InlineKeyboardButton
	for _, list := range lists {
		sListID := strconv.Itoa(*list.ID)
		callBack := createCallBackFewParam(prefixCallBackGroupProductList, sListID)
		button = tg.InlineKeyboardButton{Text: "📜 " + *list.Name, CallbackData: callBack}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tg.InlineKeyboardButton{button})
	}

	row := createInlineGroupActions(sGroupID)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row...)
	return &keyboard
}

func createInlineGroupActions(sGroupID string) [][]tg.InlineKeyboardButton {
	rows := [][]tg.InlineKeyboardButton{
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(
				"📚 New-list", *createCallBackOneParam(prefixCreateGroupList, sGroupID)),
			tg.NewInlineKeyboardButtonData(
				choiceGetAllUsersGroup, *createCallBackOneParam(prefixGetAllUsersGroup, sGroupID)),
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
