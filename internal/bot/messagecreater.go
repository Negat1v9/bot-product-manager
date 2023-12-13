package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

func createMessageProductList(p []store.Product) string {
	s := make([]byte, 0)
	row := "🍩      <b><u>Your product list</u></b>       🍪\n"
	s = append(s, row...)
	for i, product := range p {
		row = strconv.Itoa(i+1) + ". "
		row += "<b>" + product.Product + "</b>" + "\n"
		s = append(s, row...)
	}
	return string(s)
}

func createMessageCompliteGroupList(p store.FoolInfoProductList, complByUser int64) string {
	s := []byte{}
	users := map[string]int{}
	s = append(s, []byte("🪩 List - <b><u>"+*p.List.Name+"</u></b>\n")...)
	var temp string
	for _, p := range p.Products {
		linkProd := createLinkOnProduct(*p.User.UserName, p.Product)
		temp := "-  " + linkProd + "\n"
		s = append(s, temp...)
		if _, ok := users[*p.User.UserName]; !ok {
			users[*p.User.UserName] = 0
		}
		users[*p.User.UserName] += 1
	}
	s = append(s, []byte("\n💡 <b>Informations</b> 💡\n\n")...)
	for key, v := range users {
		temp = "🔸 <b>" + key + "</b>" + " - " + strconv.Itoa(v) + " added products\n"
		s = append(s, temp...)
		delete(users, key)
	}
	temp = time.Now().Format(time.DateTime)
	s = append(s, "\n✅ <u>Completed in</u> 📅\n<b>"+temp+"</b>"...)
	return string(s)
}

func createMessageComliteUserList(prod []store.InfoProduct, listName string) string {
	res := []byte("List - <b><u>" + listName + "</u></b>\n")
	t := ""
	for _, p := range prod {
		t = "-  " + "<b>" + p.Product + "</b>\n"
		res = append(res, t...)
	}
	t = time.Now().Format(time.DateTime)
	res = append(res, "\n✅ <u>Completed in</u> 📅\n<b>"+t+"</b>"...)
	return string(res)
}

func createMessageSuccessAddedProduct(p []store.Product) string {
	s := []byte{}
	row := "Success Added ✅\n"
	s = append(s, row...)
	for _, v := range p {
		row = v.Product + "\n"
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

func createLinkOnProduct(userName, product string) string {
	s := `<a href="https://t.me/` + userName + `">` + product + `</a>`
	return s
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
