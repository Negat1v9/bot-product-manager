package telegram

import (
	"fmt"
	"strconv"
)

func createMessageProductList(p []string) string {
	s := make([]byte, 0)
	row := "Youre list\n"
	s = append(s, row...)
	for i, product := range p {
		row = strconv.Itoa(i+1) + " "
		row += product + "\n"
		s = append(s, row...)
	}
	return string(s)
}

func createMessageSuccessAddedProduct(p []string) string {
	s := []byte{}
	row := "Success Added:\n"
	s = append(s, row...)
	for _, v := range p {
		row = v + "\n"
		s = append(s, row...)
	}
	return string(s)
}

func createMessgeToInviteNewUser(ownerName, groupName string) string {
	s := fmt.Sprintf(inviteUserMessage, groupName)
	return s
}
