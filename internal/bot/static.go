package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var botCommands []tgbotapi.BotCommand = []tgbotapi.BotCommand{
	{
		Command:     "/start",
		Description: "start",
	},
	{
		Command:     "/lists",
		Description: "receive all lists",
	},
	{
		Command:     "/help",
		Description: "get info",
	},
}
var (
	buttonCreateList = "new-list"
	buttonNewGroup   = "create-group"
)

// prefix for callback initilization
var (
	prefixGetUserList             = "getlist"
	prefixGetGroupLists           = "getGroups"
	prefixCallBackListProduct     = "IDNAME"
	prefixAddProductList          = "add-"
	prefixCompliteList            = "compl-"
	prefixChangeList              = "change-"
	prefixCreateGroupList         = "createGroupList-"
	prefixAddUserGroup            = "addUserGroup-"
	prefixGetUserToDelete         = "GetDeleteUser-"
	prefixCallBackListGroup       = "IDGROUP"
	prefixCallBackDelUserFromGr   = "DelUsIDGrID"
	prefixCallBackInsertUserGroup = "insertInGroup"
	prefixCallBackRefuseUserGroup = "refuseInGroup"
	prefixToMergeListGroup        = "toMerge-"
	prefixMergeListWithGroup      = "mergeListGrIDListID-"
)

// forward messages drafts
var (
	createGroupMessage       = "forward this message with the name of your group and it will be created"
	answerCreateListMsg      = "To create a new sheet, simply reply to this message with the name of the list."
	addNewProductMessage     = "Reply to this message to add new products to the list - "
	answerCreateGroupListMsg = `Reply to this message to create a list for the group -`
	answerEditListMessage    = `Reply to this message with nums products, what you want delete from - `
	textForInvitingNewUser   = `forward this message with the name of the user you want to invite - `
)

// Messages
var (
	choiceUserList        = "your lists"
	choiceGroupList       = "group lists"
	choiceWhatTypeListMsg = "Select what type of lists you want to receive\n\nYou can:\nSelect lists created for you for your self ☞ 'youre list'\n\nSelect lists created in the group ☞ 'group lists'"
	cmdStart              = "Hi friend, I'm a bot that is designed to create lists 📋 and execute them.\n\nTo find out more click /help. 💭️\n\nTo create your first list, click the ⏩ new-list ⏪ button"
	cmdHelpMessage        = "Hi, friend. Let me tell you a little about myself:\nI'm a bot that was made to automate the creation of lists in a telegram, you can:\n\n1. Create personal lists and add things to them.🔥\n\n2. Create group lists that can be edited by all its participants. 🌚"

	refuseJoinGroupMessage    = "It's a shame, but oh well, keep creating lists alone :("
	groupIsCreatesMessage     = "new group is created"
	emptyListMessage          = "now, youre list is empty"
	emptyUserInGroup          = "Youre group is empty"
	editedProductList         = "List has been success edited"
	isCompletesProductListMsg = "Congratulations, you have completed the worksheet - "

	successDeletedUser        = "User success deleted"
	refusedUserMessage        = "Unfortunately, user %s did not agree to join the group. It's better not to invite him again, why should we bother this guy in vain?"
	inviteUserMessage         = "User %s invited you to group %s, do you want to join it?"
	inviteSendMessage         = "Invitation sent"
	joinNewUserAtBotMessage   = "💥Youre friend sent you this message because he wants to invite you to a group where you can create lists together and carry them out.\nGo @golang_home_prj_bot to take advantage of it.❤️‍🔥"
	userInvitedInGroupMessage = "Congratulations, you joined the group"
	choiceWhatGroupMerge      = "Select the group you would like to connect your list with"
	successMergeListGroupMgs  = "✅ The list was successfully added to the group."
	errorMessage              = "Sorry, something seems to have gone wrong. Try later. "
)
