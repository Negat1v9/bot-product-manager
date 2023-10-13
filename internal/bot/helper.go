package telegram

import (
	"strconv"
	"strings"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

var (
	prefixCallBackListProduct   = "IDNAME"
	prefixCallBackListGroup     = "IDGROUP"
	prefixCallBackDelUserFromGr = "DelUsIDGrID"
)

func createCallBackListProducts(id int, name string) *string {
	s := prefixCallBackListProduct + strconv.Itoa(id) + name
	return &s
}

func createCallBackGroupLists(groupID int) *string {
	s := prefixCallBackListGroup + strconv.Itoa(groupID)
	return &s
}
func createCallBackDeleteUserGroup(userID, groupID string) *string {
	s := prefixCallBackDelUserFromGr + userID + "-" + groupID
	return &s
}

func parseGroupID(s string) int {
	var id int
	for i := range s {
		// isDigit
		if s[i] >= 48 && s[i] <= 57 {
			id, _ = strconv.Atoi(s[i:])
			break
		}
	}
	return id
}

// First parametr - userID second - groupID
func parseCallBackDeleteUser(s string) (uID int64, gID int) {
	sID := []byte{}
	for i := len(prefixCallBackDelUserFromGr); i < len(s); i++ {
		if s[i] == '-' {
			uID, _ = strconv.ParseInt(string(sID), 10, 64)
			sID = []byte{}
			continue
		}
		sID = append(sID, s[i])
	}
	gID, _ = strconv.Atoi(string(sID))
	return
}

func parseIDName(s string) (int, string) {
	sID := []byte{}
	var name string
	for i := len(prefixCallBackListProduct); i < len(s); i++ {
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
	for i, v := range s {
		if v == '-' && !startName {
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

func parseNameListForAddProd(s string) string {
	res, _ := strings.CutPrefix(s, addNewProductMessage)
	return res
}
func parseGroupListName(s string) string {
	res, _ := strings.CutPrefix(s, answerCreateGroupListMsg)
	return res
}
