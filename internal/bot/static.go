package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var botCommands []tgbotapi.BotCommand = []tgbotapi.BotCommand{
	{
		Command:     "/start",
		Description: "start",
	},
	{
		Command:     "/menu",
		Description: "receive all lists",
	},
	{
		Command:     "/help",
		Description: "get info",
	},
}

// prefix for callback initilization
var (
	prefixCreateSoloList           = "create_list"
	prefixCallBackListProduct      = "IDNAME@"
	prefixCallBackGroupProductList = "GRIDNAME@"
	prefixAddProductList           = "add@"
	prefixAddProductGroup          = "addGroup@"
	prefixGetUserList              = "getlist"
	prefixGetGroupLists            = "getGroups"
	prefixWantCompliteGrList       = "wantComplite@"
	prefixCompliteList             = "compl@"
	prefixCompliteSoloList         = "compliteSolo@"
	prefixSaveAsTemplete           = "saveTemplete@"
	prefixDeleteProd               = "deletProd@"
	prefixDeleteGrProd             = "deleteGrProd@"
	prefixToMergeListGroup         = "toMerge@"
	prefixMergeListWithGroup       = "mergeListGrIDListID@"
	prefixCreateGroup              = "create_group"
	prefixGetAllUsersGroup         = "GetUserGroup@"
	prefixAddUserGroup             = "addUserGroup@"
	prefixCallBackInsertUserGroup  = "insertInGroup@"
	prefixCallBackRefuseUserGroup  = "refuseInGroup@"
	prefixGetUserToDelete          = "GetDeleteUser@"
	prefixCallBackDelUserFromGr    = "DelUsIDGrID@"
	prefixCreateGroupList          = "createGroupList@"
	prefixCallBackListGroup        = "IDGROUP@"
	prefixGetMainMenu              = "getMenu"
	prefixLeaveGroup               = "leaveGroup@"
	prefixLeaveOwnerGroup          = "ownerLeaveGroupID@"
	prefixRestoreList              = "restoreList@"
	prefixRestoreGroupList         = "restoreGroupList@"
	prefixGetPageProdDelete        = "getPageProdDel@"
	prefixGetPageGroupProdDelete   = "getGrPageProdDel@"
)

// map with prefix to get type of callBack update
var prefixsMap = map[string]int{
	prefixCreateSoloList:           isWantCreateList,
	prefixCallBackListProduct:      isGetProductList,
	prefixCallBackGroupProductList: isGetGroupProductList,
	prefixAddProductList:           isWantAddNewProduct,
	prefixAddProductGroup:          isWantAddProductGroupList,
	prefixGetUserList:              isGetLists,
	prefixGetGroupLists:            isGetGroupLists,
	prefixCompliteList:             isCompliteGroupList,
	prefixCompliteSoloList:         isCompliteSoloList,
	prefixDeleteProd:               isDeleteProd,
	prefixDeleteGrProd:             isDeleteGroupProd,
	prefixToMergeListGroup:         isWantMergeList,
	prefixMergeListWithGroup:       isMergeListGroup,
	prefixCreateGroup:              isWantCreateGroup,
	prefixGetAllUsersGroup:         isGetAllUsersGroup,
	prefixAddUserGroup:             isWantInviteNewUser,
	prefixCallBackInsertUserGroup:  isUserReadyJoinGroup,
	prefixCallBackRefuseUserGroup:  isUserRefusedGroup,
	prefixGetUserToDelete:          isGetUsersToDelete,
	prefixCallBackDelUserFromGr:    isDeleteUserFromGroup,
	prefixCreateGroupList:          isWantCreateGroupList,
	prefixCallBackListGroup:        isGetAllGroupLists,
	prefixGetMainMenu:              isGetMainMenu,
	prefixLeaveGroup:               isLeaveGroup,
	prefixLeaveOwnerGroup:          isLeaveOwnerGroup,
	prefixWantCompliteGrList:       isWantCompliteGrList,
	prefixRestoreList:              isRestoreProductList,
	prefixRestoreGroupList:         isRestoreGroupList,
	prefixGetPageProdDelete:        isWantDeleteProductFromList,
	prefixGetPageGroupProdDelete:   isWantDeleteProductFromGroupList,
}

// forward messages drafts
var (
	createGroupMessage        = "🧐 What band name would you prefer for group?"
	answerCreateListMsg       = "✏️ Write a <b>name</b> for the new List"
	addNewProductMessageReply = "Write via 👉 <b>.</b> or 👉 <b>,</b> product names to add it in List"
	answerCreateGroupListMsg  = "✒️ Write a name for the list - "
	answerEditListMessage     = `🎱 Write the <u>numbers</u> of the products you would like to remove from the list - `
	textForInvitingNewUser    = `Write the user's <b>@nickname</b> to send him an invitation to your group - `
)

// Messages
var (
	choiceUserList  = "📜 My List"
	choiceGroupList = "👥 My groups"

	choiceGetAllUsersGroup = "Users 🧞‍♂️"
	cmdMenu                = "🗿 <b>Options</b> 🗿\n\n🎱 Push <b><u>My List</u></b> button to get all your lists or create a new one\n\n🏀 Push <b><u>My Group</u></b> to get all your group or create a new one"
	cmdStart               = "Hi friend, I'm a bot that is designed to create lists 📋 and execute them.\n\nTo find out more click /help. 💭️\n\nClick on /menu to receive all options"
	cmdHelpMessage         = "Hi, friend 👋. Let me tell you a little about myself:\n\nI'm a bot 👾 that was made to automate\nthe creation of lists 📝 in a telegram, you can:\n\n1. Create personal lists and add things to them.🔥\n\n2. Create group lists that can be edited by all its participants. 🌚\n\n❓ How to use ❓\n\n1️⃣ Select the list to add a new product\n\n2️⃣ Click the add button and\n\n3️⃣ Write message with the product\nnames separated by a ',' or '.'\n\n🟠Example 🧾\n\n✏️ Squash caviar, Juice, Potato, Soup ✏️"

	refuseJoinGroupMessage    = "It's a shame, but oh well, keep creating lists alone :("
	groupIsCreatesMessage     = "New group is created"
	emptyListMessage          = "😿 Now, youre list is empty"
	emptyUserInGroup          = "Youre group is empty"
	editedProductList         = "List has been success edited"
	isCompletesProductListMsg = "Congratulations, you have completed the worksheet ✅"

	ownerGroupWantLeave       = "❗ You are the owner of this group, if you leave it it will be deleted ❌"
	successLeaveGroup         = "You have successfully left the group 🌈"
	successDeletedUser        = "User success deleted ❎"
	refusedUserMessage        = "Unfortunately, user %s did not agree to join the group. It's better not to invite him again, why should we bother this guy in vain?"
	inviteUserMessage         = "User %s invited you to group %s, do you want to join it?"
	inviteSendMessage         = "Invitation sent"
	joinNewUserAtBotMessage   = "💥 Youre friend sent you this message because he wants to invite you to a group where you can create lists together and carry them out.\nGo @golang_home_prj_bot to take advantage of it.❤️‍🔥"
	userInvitedInGroupMessage = "Congratulations, you joined the group"
	choiceWhatGroupMerge      = "Select the group you would like to connect your list with"
	successMergeListGroupMgs  = "✅ The list was successfully added to the group."
	listsProductsMsgHelp      = "Click on your list name to go to it 👇"
	NotNickNameUserMsg        = "😥 Unfortunately, I can`t create a group or invite you to other groups if you don`t have a <b>NickName</b>.\n\n💥But you can create it!\n\n<u>Information on how to do this</u> <a href=\"https://screenrant.com/create-change-telegram-username-how/#:~:text=Set%20Your%20Telegram%20Username,it%20create%20a%20unique%20username.\">Here</a>\n\n💭 If you have created a <b>Nickname</b> for yourself click /start,\nand you will be able to create groups and you will be able\nto be invited to third-party groups"
	errorMessage              = "Sorry, something seems to have gone wrong. Try later. "
	errorLastCmdUserDeleted   = "🫣 Unfortunately, I'm very forgetful bot 🤖, if you wanted to do something, please try again 🙏"
)
