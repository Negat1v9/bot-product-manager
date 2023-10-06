package telegram

import (
	"strings"
)

func isCreateList(s string) bool {
	return strings.EqualFold(s, buttonCreateList)
}

func isSelectUserList(s string) bool {
	return strings.EqualFold(s, buttonListList)
}

func isCreateNameForward(s string) bool {
	return strings.EqualFold(s, answerCreateListMsg)
}

func isAddNewProductForward(s string) bool {
	return strings.HasPrefix(s, addNewProductMessage)
}

func isGetProductList(s string) bool {
	return strings.HasPrefix(s, "IDNAME")
}

func isAddNewProduct(s string) bool {
	return strings.HasPrefix(s, "add")
}
func isCreateAddProduct(s string) bool {
	return strings.HasPrefix(s, "create-")
}

// // TODO: make faster
// func isAddCommand(s string) bool {
// 	return strings.ContainsAny(s, button)
// }
