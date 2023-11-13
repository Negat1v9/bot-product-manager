package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// send user message to reply with new list
	isWantCreateList = iota
	// is User replyed with name
	isCreateNewList // message
	// is get list with products
	isGetProductList
	// is get group list
	isGetGroupProductList

	// send user message to reply with new products
	isWantAddNewProduct
	// user replyed on message with new products
	isAddNewProduct // message
	// get all lists names
	isGetLists
	// get all groups names
	isGetGroupLists
	// get in group all templates what was saved
	isGetGroupTemplates
	// get one group template for edit
	isGetOneGroupTemplate
	// user want complite group list and he receive choice save template
	isWantCompliteList
	// complite list
	isCompliteList
	// complite solo list without chance save template
	isCompliteSoloList
	// user want choice save list as template
	isSaveTemplete
	// send user message to reply with idited nums pruducts
	isWantEditList
	// user replyed on message with nums products to edit
	isEditList // message
	// send message with groups name to merge list
	isWantMergeList
	// user select groups name to merge list
	isMergeListGroup
	// user after creation list choice merge with template
	isWantConnectTemplate
	// get template list for merging with new list
	isGetTemplateForConnect
	// user choice is connect Merge new list with template
	isConnectTemplate

	// is send message to reply with new groups name
	isWantCreateGroup
	// user replyed with new group name
	isCreateGroup // message
	// send all users from group
	isGetAllUsersGroup
	// is send message to reply with new userName
	isWantInviteNewUser
	// user replyed on message with new userName
	isSendInviteNewUser // message text
	// is user ready to join in new group
	isUserReadyJoinGroup
	// is user refused on invite
	isUserRefusedGroup
	// is send all users to delete from group
	isGetUsersToDelete
	// is user chouce who will deleted from group
	isDeleteUserFromGroup
	// send message to reply with new group list name
	isWantCreateGroupList
	// user replyed on message with new group list name
	isCreateGroupList // message text
	// is send all groups lists name
	isGetAllGroupLists
	// is get menu with all options
	isGetMainMenu
	// user want leave from group
	isLeaveGroup
	// is Owner group leave the group and it deleted
	isLeaveOwnerGroup
)

// INFO: CustomType for message for log time respose create
type MessageWithTime struct {
	Msg        *tgbotapi.MessageConfig
	EditMesage *tgbotapi.EditMessageTextConfig
	WorkTime   time.Time
}
