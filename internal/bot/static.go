package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var botCommands []tgbotapi.BotCommand = []tgbotapi.BotCommand{
	{
		Command:     "/start",
		Description: "start",
	},
	{
		Command:     "/help",
		Description: "get info",
	},
}
var (
	buttonListList     = "my-list"
	buttonCreateList   = "new-list"
	buttonGetUserGroup = "get-group"
	buttonNewGroup     = "create-group"
)
var commandStart = ""
var createGroupMessage string = "Перешлите это сообщение, с названием вашей группы :)"
var groupIsCreatesMessage string = "Новая группа создана"
var emptyListMessage string = "Пока что ваш список пуст"
var emptyUserInGroup string = "Youre group is empty"
var cmdHelpMessage string = `Этот бот предназначен для создания удобных списков продуктов или вещей, для создания групп, где люди смогут закрывать весь список задач.`
var isCompletesProductListMsg string = "Congratulations, you have completed the worksheet - "
var answerCreateGroupListMsg = `Ответе на это сообщение, чтобы создать список для группы -`
var answerCreateListMsg string = `Чтобы создать новый лист, просто ответе на это сообщение с названием списка.`
var addNewProductMessage string = `Ответе на это сообщение, чтобы добавить новые продукты в список - `
var successDeletedUser string = `User success deleted`
var textForInvitingNewUser string = `forward this message with the name of the user you want to invite - `
var inviteUserMessage string = `User %s invited you to group %s, do you want to join it?`
var inviteSendMessage string = `Invitation sent`
var userInvitedInGroupMessage string = "Congratulations, you joined the group %s"
