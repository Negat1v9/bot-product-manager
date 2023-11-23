package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

func createMessageProductList(p []string) string {
	s := make([]byte, 0)
	row := "Youre product list 📋\n"
	s = append(s, row...)
	for i, product := range p {
		row = strconv.Itoa(i+1) + ". "
		row += "<b>" + product + "</b>" + "\n"
		s = append(s, row...)
	}
	return string(s)
}

func createMessageCompliteList(p store.Product, complByUser int64) string {
	s := []byte{}
	s = append(s, []byte("🪩 List - <b><u>"+p.ListName+"</u></b>\n")...)
	var temp string
	for _, product := range p.Products {
		temp = "-  "
		temp += "<b>" + product + "</b>" + "\n"
		s = append(s, temp...)
	}
	s = append(s, []byte("\n💡 <b>Informations</b> 💡\n")...)
	for _, editor := range p.Editors {
		temp = "🔸 <b>" + *editor.User.UserName + "</b>" + " - " + strconv.Itoa(editor.ManyAddProducts) + " added products\n"
		s = append(s, temp...)
	}
	temp = time.Now().Format(time.DateTime)
	s = append(s, []byte("\n✅ <u>Completed in</u> 📅\n<b>"+temp+"</b>")...)
	return string(s)
}

func createMessageSuccessAddedProduct(p []string) string {
	s := []byte{}
	row := "Success Added ✅\n"
	s = append(s, row...)
	for _, v := range p {
		row = v + "\n"
		s = append(s, row...)
	}
	return string(s)
}

func createMessageGetAllUsersGroup(users []store.User, ownerID int64) string {
	r := []byte("People near you:\n")
	s := ""
	for i, user := range users {
		if user.ChatID != ownerID {
			s = strconv.Itoa(i+1) + ". " + *user.UserName + "\n"
		} else {
			s = strconv.Itoa(i+1) + ". " + *user.UserName + " 👑\n"
		}
		r = append(r, s...)
	}
	return string(r)
}

func createButtonDeleteUser(userName string) string {
	return "❌ " + userName + " ❌"
}

func createMessageUserRefusedOrder(refusedName string) string {
	s := fmt.Sprintf(refusedUserMessage, refusedName)
	return s
}
func createMessgeToInviteNewUser(ownerName, groupName string) string {
	s := fmt.Sprintf(inviteUserMessage, ownerName, groupName)
	return s
}
