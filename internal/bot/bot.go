package telegram

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	manager "github.com/Negat1v9/telegram-bot-orders/internal"
	"github.com/Negat1v9/telegram-bot-orders/store/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

const (
	bufferMessages       = 10
	logDuractionResponse = "time create and send response"
)

type Bot struct {
	timeOut int
	offset  int
	output  chan *MessageWithTime
	client  *tgbotapi.BotAPI
	logger  *slog.Logger
	hub     manager.Manager
}

func New(client *tgbotapi.BotAPI, timeOut int, offset int) *Bot {
	return &Bot{
		timeOut: timeOut,
		offset:  offset,
		output:  make(chan *MessageWithTime, bufferMessages),
		client:  client,
		logger:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (b *Bot) Start(dbURL string) error {

	db, err := initDB(dbURL)
	if err != nil {
		return err
	}

	store := sqlite.Newstorage(db)
	// create all tables ONCE
	err = store.CreateTables()
	if err != nil {
		return err
	}

	b.hub = NewHub(store, b.output)

	configBot := tgbotapi.NewUpdate(b.offset)

	configBot.Timeout = b.timeOut
	// channel for getting updates
	updates := b.client.GetUpdatesChan(configBot)
	// set commands for bot
	b.client.Send(tgbotapi.NewSetMyCommands(botCommands...))
	go b.MessageSender()
	return b.startPooling(updates)
}

func (b *Bot) startPooling(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		// Message update
		if update.Message != nil {
			// manage cmd or message to hub
			go func(msg *tgbotapi.Message) {
				start := time.Now()
				err := b.hub.MessageUpdate(msg, start)

				if err != nil {
					b.logger.Error("Manage - Create Message:",
						slog.String("error", err.Error()))
				}
				return
				// if err = b.sendMessage(res); err != nil {
				// 	b.logger.Error("not sended message",
				// 		slog.String("error", err.Error()))
				// }
				// b.logger.Info("time message routing", slog.Duration("time", time.Now().Sub(s)))
			}(update.Message)
			continue
		}
		if update.CallbackQuery != nil {
			go func(c *tgbotapi.CallbackQuery) {
				start := time.Now()
				err := b.hub.CallBackUpdate(c, start)

				if err != nil {
					b.logger.Error("Manage callback update:",
						slog.String("error", err.Error()))
				}

				// if err = b.sendMessage(res); err != nil {

				// }
				// b.logger.Info("time message routing", slog.Duration("time", time.Now().Sub(s)))
			}(update.CallbackQuery)
		}
		continue
	}
	return nil
}

func (b *Bot) MessageSender() error {
	for {
		select {
		case msg := <-b.output:
			go func() {
				defer func() {
					b.logger.Info(logDuractionResponse,
						slog.Duration(" ", time.Now().Sub(msg.WorkTime)))
				}()
				if _, err := b.client.Send(msg.Msg); err != nil {
					b.logger.Error("not sended message",
						slog.String("error", err.Error()))
					return
				}
			}()
		default:
			continue
		}
	}
}
func initDB(URL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", URL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
