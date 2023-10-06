package telegram

import (
	"strconv"
	"strings"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

var (
	prefixCallBackProductList = "IDNAME"
)

func createCallBackListProducts(id int, name string) *string {
	s := prefixCallBackProductList + strconv.Itoa(id) + name
	return &s
}

func parseIDNameList(s string) (int, string) {
	sID := []byte{}
	var name string
	for i := len(prefixCallBackProductList); i < len(s); i++ {
		if s[i] >= 48 && s[i] <= 57 {
			sID = append(sID, s[i])
		} else {
			name = s[i:]
			break
		}
	}
	id, _ := strconv.Atoi(string(sID))
	return id, name
}

// For example "add-122 || comlp-22"
func parseNameListFromProductAction(s string) string {
	name := []byte{}
	startName := false
	for i := range s {
		if string(s[i]) == "-" && !startName {
			startName = true
			continue
		}
		if startName {
			name = append(name, s[i])
		}
	}
	return string(name)
}

// pr, pr
// Parse Latter is: "," || "."
func parseStringToProducts(s string, listID int) store.Product {
	products := []string{}
	prev := 0
	for i := range s {
		if s[i] == ',' || s[i] == '.' {
			products = append(products, cutLineBreak(s[prev:i]))
			prev = i + 1
		}
	}
	// if String does not ends on "." or ","
	if prev != len(s) {
		products = append(products, cutLineBreak(s[prev:]))
	}
	res := store.Product{
		Products: products,
		ListID:   listID,
	}
	return res
}
func cutLineBreak(s string) string {
	l, r := 0, len(s)-1
	for i := range s {
		if s[i] == '\n' || s[i] == ' ' {
			l++
		} else if s[r] == '\n' || s[r] == ' ' {
			r--
		} else {
			break
		}
	}
	return s[l : r+1]
}

func parseNameList(s string) string {
	res, _ := strings.CutPrefix(s, addNewProductMessage)
	return res
}
