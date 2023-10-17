package telegram

import (
	"strings"
)

func isCommandUpdate(u string) bool {
	if string(u[0]) == "/" {
		return true
	}
	return false
}

// Product List filters
func isCreateNameForward(s string) bool {
	return strings.EqualFold(s, answerCreateListMsg)
}

func isCreateList(s string) bool {
	return strings.EqualFold(s, buttonCreateList)
}

func isGetProductList(s string) bool {
	return strings.HasPrefix(s, prefixCallBackListProduct)
}

func isSelectUserList(s string) bool {
	return strings.EqualFold(s, buttonListList)
}

// PRODUCTS FILTERS

func isAddNewProductForward(s string) bool {
	return strings.HasPrefix(s, addNewProductMessage)
}

func isAddNewProduct(s string) bool {
	return strings.HasPrefix(s, prefixAddProductList)
}

func isCompliteProductList(s string) bool {
	return strings.HasPrefix(s, prefixCompliteList)
}

func isEditProductList(s string) bool {
	return strings.HasPrefix(s, prefixChangeList)
}

func isEditListForward(s string) bool {
	return strings.HasPrefix(s, answerEditListMessage)
}

// GROUP FILTERS
func isCreateGroup(s string) bool {
	return strings.EqualFold(s, buttonNewGroup)
}

func isAddNewUserGroup(s string) bool {
	return strings.HasPrefix(s, prefixAddUserGroup)
}

func isGetUsersForDelGroup(s string) bool {
	return strings.HasPrefix(s, prefixGetUserToDelete)
}

func isGetUserGroup(s string) bool {
	return strings.EqualFold(s, buttonGetUserGroup)
}

func isCreateGroupList(s string) bool {
	return strings.HasPrefix(s, prefixCreateGroupList)
}

func isCreateGroupListForward(s string) bool {
	return strings.HasPrefix(s, answerCreateGroupListMsg)
}

func isCreateNewGroupForward(s string) bool {
	return strings.EqualFold(s, createGroupMessage)
}

func isGetGroupLists(s string) bool {
	return strings.HasPrefix(s, prefixCallBackListGroup)
}

func isDeleteUserFromGroup(s string) bool {
	return strings.HasPrefix(s, prefixCallBackDelUserFromGr)
}

func isUserReadyInvite(s string) bool {
	return strings.HasPrefix(s, prefixCallBackInsertUserGroup)
}

func isUserRefuseInvite(s string) bool {
	return strings.HasPrefix(s, prefixCallBackRefuseUserGroup)
}

func isSendInviteToNewUser(s string) bool {
	return strings.HasPrefix(s, textForInvitingNewUser)
}
