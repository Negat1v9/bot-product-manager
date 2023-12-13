package telegram

import (
	// "strconv"
	// "fmt"
	"strconv"
	"strings"

	helper "github.com/Negat1v9/telegram-bot-orders/internal/helpers"
	"github.com/Negat1v9/telegram-bot-orders/store"
)

// create callback with one parametr
func createCallBackOneParam(prefix string, data string) *string {
	res := prefix + data
	return &res
}

func parseCallBackOneParam(prefix, callBack string) string {
	return callBack[len(prefix):]
}

// create callback view like prefisDATA1@DATA2...
func createCallBackFewParam(prefix string, data ...string) *string {
	res := []byte{}
	res = append(res, []byte(prefix+data[0])...)
	for i := 1; i < len(data); i++ {
		res = append(res, []byte("@"+data[i])...)
	}
	s := string(res)
	return &s
}
func parseCallBackFewParam(prefix, callBack string) []string {
	res := []string{}
	start := len(prefix)
	for i := len(prefix); i < len(callBack); i++ {
		if callBack[i] == '@' {
			res = append(res, callBack[start:i])
			start = i + 1
			continue
		}
	}
	res = append(res, callBack[start:])
	return res
}

// convert string to int second parametr is typeInt default return 0
func convSToI[T int | int64](p string, base int) T {
	switch base {
	case 0:
		res, _ := strconv.ParseInt(p, 10, 64)
		return T(res)
	case 64:
		res, _ := strconv.ParseInt(p, 10, 0)
		return T(int(res))
	default:
		return 0
	}
}

// Parse Latter is: "," || "."
func parseStringToProducts(s string, ownerID int64, listID int) []store.Product {
	products := []store.Product{}
	prev := 0
	for i := range s {
		if s[i] == ',' || s[i] == '.' {
			prod := store.Product{
				Product: helper.CutLineBreak((s[prev:i])),
				UserID:  ownerID,
				ListID:  &listID,
			}
			products = append(products, prod)
			prev = i + 1
		}
	}
	// if String does not ends on "." or ","
	if prev != len(s) {
		prod := store.Product{
			Product: helper.CutLineBreak((s[prev:])),
			UserID:  ownerID,
			ListID:  &listID,
		}
		products = append(products, prod)
	}

	return products
}

func parseIndexEditProduct(s string) map[int]bool {
	res := make(map[int]bool)
	temp := ""
	for i := range s {
		if s[i] >= 48 && s[i] <= 57 {
			temp += string(s[i])

		} else if temp != "" {
			t, err := strconv.Atoi(temp)
			if err != nil {
				temp = ""
				continue
			}
			res[t-1] = true
			temp = ""
		}
	}
	if temp != "" {
		t, _ := strconv.Atoi(temp)
		res[t-1] = true
	}
	return res
}

func deleteProductByIndex(products []string, targets map[int]bool) []string {
	capas := len(products) - len(targets)
	if capas <= 0 {
		capas = len(products)
	}
	res := make([]string, 0, capas)
	for i := range products {
		// if exist in target map -> delete
		if _, ok := targets[i]; !ok {
			res = append(res, products[i])
		}
	}
	return res
}

func makeNameClear(name string) string {
	name = helper.CutLineBreak(name)
	res := make([]rune, 0, len(name))
	for _, v := range name {
		if v == ' ' {
			res = append(res, '_')
		} else {
			res = append(res, v)
		}
	}
	return string(res)
}

func clearCallBackData(s string) string {
	for i, v := range s {
		if v == '@' {
			return s[:i+1]
		}
	}
	return s
}

// example NickName "@NaemNuam"
func parseUserNickNameForAddGroup(s string) string {
	return s[1:]
}

// func parseNameGroupAddUser(s string) string {
// 	res, _ := strings.CutPrefix(s, textForInvitingNewUser)
// 	return res
// }

// func parseNameListActions(s string) string {
// 	for i := len(s) - 1; i >= 0; i-- {
// 		if s[i] == ' ' {
// 			return s[i+1:]
// 		}
// 	}

// 	return " "
// }

// Check slice of users and compare with target ID
func checkUserInGroup(targetUserID int64, users []store.User) bool {
	for _, user := range users {
		if user.ChatID == targetUserID {
			return true
		}
	}
	return false
}

func searchOwnerGroup(ownerID int64, users []store.User) *store.User {
	for _, user := range users {
		if user.ChatID == ownerID {
			return &user
		}
	}
	return nil
}

// Info: add many of new product in Editors for current list, if its first edits, create one
// func addManyEditsProductList(u store.User, editors []store.Editors, manyNewProducts int) []store.Editors {
// 	isExistUser := false
// 	indexUser := 0
// 	for i, e := range editors {
// 		if e.User.ChatID == u.ChatID {
// 			indexUser = i
// 			isExistUser = true
// 			break
// 		}
// 	}
// 	if isExistUser {
// 		editors[indexUser].ManyAddProducts += manyNewProducts
// 	} else {
// 		editor := store.Editors{User: u, ManyAddProducts: manyNewProducts}
// 		editors = append(editors, editor)
// 	}
// 	return editors
// }

// Input first row from splited text bu '\n'
func parseNameTextList(splitedText []string) string {
	if len(splitedText) == 0 {
		return " "
	}
	fRow := splitedText[0]
	for i := len(fRow) - 1; i >= 0; i-- {
		if fRow[i] == ' ' {
			return fRow[i+1:]
		}
	}
	return " "
}

func parseTextToProd(text []string, ownerID int64, listID int) []store.Product {
	prod := []store.Product{}
	for _, v := range text {
		if v[0] == '-' {
			t, _ := strings.CutPrefix(v, "-  ")
			p := store.Product{
				Product: t,
				UserID:  ownerID,
				ListID:  &listID,
			}
			prod = append(prod, p)
		}
	}
	return prod
}

func parseTextToProdGroup(text []string, ownerID int64, listID int) []store.Product {
	prod := []store.Product{}
	for _, v := range text {
		if v[0] == '-' {
			p := store.Product{
				Product: cutLinkOnUser(v),
				UserID:  ownerID,
				ListID:  &listID,
			}
			prod = append(prod, p)
		}
	}
	return prod
}
func cutLinkOnUser(row string) string {
	r := len(row) - 4
	for i := r; i >= 0; i-- {
		if row[i] == '>' {
			return row[i+1 : r]
		}
	}
	return ""
}
func splitText(text string, key rune) []string {
	r := []string{}
	t := ""
	for _, v := range text {
		if v == key {
			if t != "" {
				r = append(r, t)
				t = ""
			}
		} else {
			t += string(v)
		}
	}
	return r
}
