package telegram

import (
	"strconv"
	"strings"

	"github.com/Negat1v9/telegram-bot-orders/store"
)

var (
	prefixCallBackListProduct     = "IDNAME"
	prefixCallBackListGroup       = "IDGROUP"
	prefixCallBackDelUserFromGr   = "DelUsIDGrID"
	prefixCallBackInsertUserGroup = "insertInGroup"
	prefixCallBackRefuseUserGroup = "refuseInGroup"
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
func createCallBackInsertNewUserGroup(userID, groupID string) *string {
	s := prefixCallBackInsertUserGroup + userID + "-" + groupID
	return &s
}

func createCallBackRefuseGroup(userID, groupID string) *string {
	s := prefixCallBackRefuseUserGroup + userID + "-" + groupID
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

// NOTE: for add and refuse inviting maybe groupID is not nessesory
// First parametr - userID second - groupID
func parseCallBackGroupActions(s string) (uID int64, gID int) {
	tempID := []byte{}
	isStartIds := false // if names contains "-"
	for i := 0; i < len(s); i++ {
		if s[i] == '-' && !isStartIds {
			uID, _ = strconv.ParseInt(string(tempID), 10, 64)
			tempID = []byte{}
			isStartIds = true
			continue
		}

		if s[i] >= 48 && s[i] <= 57 {
			tempID = append(tempID, s[i])
		}
	}
	gID, _ = strconv.Atoi(string(tempID))
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
