package telegram

func isCommandUpdate(u string) bool {
	if string(u[0]) == "/" {
		return true
	}
	return false
}

// Product List filters

// func isCreateNameForward(s string) bool {
// 	return strings.EqualFold(s, answerCreateListMsg)
// }

// // // PRODUCTS FILTERS

// func isAddNewProductForward(s string) bool {
// 	return strings.Contains(s, addNewProductMessageReply)
// }

// func isAddNewProductGroupForward(s string) bool {
// 	return strings.Contains(s, addNewProductAtGroupList)

// }

// func isEditListForward(s string) bool {
// 	return strings.Contains(s, answerEditListMessage)
// }
// func isEditGroupListForward(s string) bool {
// 	return strings.Contains(s, answerEditGroupList)
// }

// func isCreateGroupListForward(s string) bool {
// 	return strings.HasPrefix(s, answerCreateGroupListMsg)
// }

// func isSendInviteToNewUser(s string) bool {
// 	return strings.HasPrefix(s, textForInvitingNewUser)
// }

// func isCreateNewGroupForward(s string) bool {
// 	return strings.EqualFold(s, createGroupMessage)
// }
