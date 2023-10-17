package telegram

import (
	// "strconv"
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

// create callback view like prefisDATA1-DATA2...
func createCallBackFewParam(prefix string, data ...string) *string {
	res := []byte{}
	res = append(res, []byte(prefix+data[0])...)
	for i := 1; i < len(data); i++ {
		res = append(res, []byte("-"+data[i])...)
	}
	s := string(res)
	return &s
}
func parseCallBackFewParam(prefix, callBack string) []string {
	res := []string{}
	start := len(prefix)
	for i := len(prefix); i < len(callBack); i++ {
		if callBack[i] == '-' {
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
func parseStringToProducts(s string, listID int) store.Product {
	products := []string{}
	prev := 0
	for i := range s {
		if s[i] == ',' || s[i] == '.' {
			products = append(products, helper.CutLineBreak((s[prev:i])))
			prev = i + 1
		}
	}
	// if String does not ends on "." or ","
	if prev != len(s) {
		products = append(products, helper.CutLineBreak(s[prev:]))
	}
	res := store.Product{
		Products: products,
		ListID:   listID,
	}
	return res
}

// NOTE: Test this
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

// example NickName "@NaemNuam"
func parseUserNickNameForAddGroup(s string) string {
	return s[1:]
}

func parseNameGroupAddUser(s string) string {
	res, _ := strings.CutPrefix(s, textForInvitingNewUser)
	return res
}

func parseNameListForAddProd(s string) string {
	res, _ := strings.CutPrefix(s, addNewProductMessage)
	return res
}
func parseGroupListName(s string) string {
	res, _ := strings.CutPrefix(s, answerCreateGroupListMsg)
	return res
}

func parseListNameEditList(s string) string {
	res, _ := strings.CutPrefix(s, answerEditListMessage)
	return res
}

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
