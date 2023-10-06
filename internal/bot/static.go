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
	buttonListList   = "my-list"
	buttonCreateList = "new-list"
	buttonNewGroup   = "new-group"
)

var emptyListMessage string = "Пока что ваш список пуст"
var cmdHelpMessage string = `Этот бот предназначен для создания удобных списков продуктов или вещей, для создания групп, где люди смогут закрывать весь список задач.`

var answerCreateListMsg string = `Чтобы создать новый лист, просто ответе на это сообщение с названием списка.`
var addNewProductMessage string = `Ответе на это сообщение, чтобы добавить новые продукты в список - `
