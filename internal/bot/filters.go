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

func isCreateList(s string) bool {
	return strings.EqualFold(s, buttonCreateList)
}

func isSelectUserList(s string) bool {
	return strings.EqualFold(s, buttonListList)
}

func isCreateGroup(s string) bool {
	return strings.EqualFold(s, buttonNewGroup)
}

func isGetUserGroup(s string) bool {
	return strings.EqualFold(s, buttonGetUserGroup)
}

func isCreateNameForward(s string) bool {
	return strings.EqualFold(s, answerCreateListMsg)
}

func isCreateGroupListForward(s string) bool {
	return strings.HasPrefix(s, answerCreateGroupListMsg)
}
func isAddNewProductForward(s string) bool {
	return strings.HasPrefix(s, addNewProductMessage)
}

func isCreateNewGroupForward(s string) bool {
	return strings.EqualFold(s, createGroupMessage)
}

func isGetProductList(s string) bool {
	return strings.HasPrefix(s, prefixCallBackListProduct)
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

func isAddNewProduct(s string) bool {
	return strings.HasPrefix(s, "add-")
}
func isCreateAddProduct(s string) bool {
	return strings.HasPrefix(s, "create-")
}
func isCreateGroupList(s string) bool {
	return strings.HasPrefix(s, "createGroupList-")
}

func isAddNewUserGroup(s string) bool {
	return strings.HasPrefix(s, "addUserGroup-")
}

func isGetUsersForDelGroup(s string) bool {
	return strings.HasPrefix(s, "GetDeleteUser-")
}
func isCompliteProductList(s string) bool {
	return strings.HasPrefix(s, "comple-")
}
