package telegram

import "strings"

func isCommandUpdate(u string) bool {
	if string(u[0]) == "/" {
		return true
	}
	return false
}

// Product List filters
//
//	func isCreateList(s string) bool {
//		return strings.EqualFold(s, prefixCreateSoloList)
//	}
func isCreateNameForward(s string) bool {
	return strings.EqualFold(s, answerCreateListMsg)
}

// func isGetProductList(s string) bool {
// 	return strings.HasPrefix(s, prefixCallBackListProduct)
// }

// // PRODUCTS FILTERS

func isAddNewProductForward(s string) bool {
	return strings.HasPrefix(s, addNewProductMessageReply)
}

// func isAddNewProduct(s string) bool {
// 	return strings.HasPrefix(s, prefixAddProductList)
// }

// func isUserChoiceLists(s string) bool {
// 	return strings.HasPrefix(s, prefixGetUserList)
// }
// func isUserChoiceGroupLists(s string) bool {
// 	return strings.HasPrefix(s, prefixGetGroupLists)
// }

// func isCompliteProductList(s string) bool {
// 	return strings.HasPrefix(s, prefixCompliteList)
// }

//	func isEditProductList(s string) bool {
//		return strings.HasPrefix(s, prefixChangeList)
//	}
func isEditListForward(s string) bool {
	return strings.HasPrefix(s, answerEditListMessage)
}

// // Info: is get group to merge list
// func isWantMergeList(s string) bool {
// 	return strings.HasPrefix(s, prefixToMergeListGroup)
// }

// // Info: is push button with name group
// func isMergeListGroup(s string) bool {
// 	return strings.HasPrefix(s, prefixMergeListWithGroup)
// }

// // GROUP FILTERS
//
//	func isCreateGroup(s string) bool {
//		return strings.EqualFold(s, prefixCreateGroup)
//	}
func isCreateGroupListForward(s string) bool {
	return strings.HasPrefix(s, answerCreateGroupListMsg)
}

// func isGetAllUsersGroup(s string) bool {
// 	return strings.HasPrefix(s, prefixGetAllUsersGroup)
// }

//	func isAddNewUserGroup(s string) bool {
//		return strings.HasPrefix(s, prefixAddUserGroup)
//	}
func isSendInviteToNewUser(s string) bool {
	return strings.HasPrefix(s, textForInvitingNewUser)
}

// func isUserReadyInvite(s string) bool {
// 	return strings.HasPrefix(s, prefixCallBackInsertUserGroup)
// }
// func isUserRefuseInvite(s string) bool {
// 	return strings.HasPrefix(s, prefixCallBackRefuseUserGroup)
// }

// func isGetUsersForDelGroup(s string) bool {
// 	return strings.HasPrefix(s, prefixGetUserToDelete)
// }

// func isDeleteUserFromGroup(s string) bool {
// 	return strings.HasPrefix(s, prefixCallBackDelUserFromGr)
// }

// func isCreateGroupList(s string) bool {
// 	return strings.HasPrefix(s, prefixCreateGroupList)
// }

func isCreateNewGroupForward(s string) bool {
	return strings.EqualFold(s, createGroupMessage)
}

// func isGetGroupLists(s string) bool {
// 	return strings.HasPrefix(s, prefixCallBackListGroup)
// }
