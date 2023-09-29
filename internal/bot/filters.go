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

// // TODO: make faster
// func isAddCommand(s string) bool {
// 	return strings.ContainsAny(s, button)
// }
