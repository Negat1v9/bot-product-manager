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
	timeForSkipUpdates   = 5
	bufferMessages       = 10
	logDuractionResponse = "time create and send response"
)

type Bot struct {
	timeOut int
	offset  int
	output  chan MessageWithTime
	client  *tgbotapi.BotAPI
	logger  *slog.Logger
	hub     manager.Manager
}

func New(client *tgbotapi.BotAPI, timeOut int, offset int) *Bot {
	return &Bot{
		timeOut: timeOut,
		offset:  offset,
		output:  make(chan MessageWithTime, bufferMessages),
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
	b.skipLastUpdates(updates)
	for update := range updates {
		// Message update
		if update.Message != nil {
			b.logger.Info("Get messageid", slog.Int("ID", update.Message.MessageID))
			// manage cmd or message to hub
			go func(msg *tgbotapi.Message) {
				start := time.Now()
				err := b.hub.MessageUpdate(msg, start)

				if err != nil {
					b.logger.Error("Manage - Create Message:",
						slog.String("error", err.Error()))
				}
				return

			}(update.Message)
			continue
		}
		if update.CallbackQuery != nil {
			b.logger.Info("Get messageid", slog.Int("ID", update.CallbackQuery.Message.MessageID))
			go func(c *tgbotapi.CallbackQuery) {
				start := time.Now()
				err := b.hub.CallBackUpdate(c, start)

				if err != nil {
					b.logger.Error("Manage callback update:",
						slog.String("error", err.Error()))
				}

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
			if msg.Msg == nil && msg.EditMesage == nil {
				b.logger.Error("response nil message, message can`t be send")
				continue
			}
			go func() {
				defer func() {
					b.logger.Info(logDuractionResponse,
						slog.Duration(" ", time.Now().Sub(msg.WorkTime)))
				}()
				if msg.EditMesage != nil {
					if _, err := b.client.Send(msg.EditMesage); err != nil {
						b.logger.Error("not sended edit message",
							slog.String("error", err.Error()))
					}
					return
				}
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
func (b *Bot) skipLastUpdates(updates tgbotapi.UpdatesChannel) {
	timer := time.NewTicker(time.Second * timeForSkipUpdates)
	end := make(chan struct{})
	go func() {
		count := 0
		for {
			select {
			case <-updates:
				b.logger.Info("skip update", slog.IntValue(count))
				count++
			case <-timer.C:
				end <- struct{}{}
				return
			}
		}
	}()
	<-end
	b.logger.Info("Skiping end")
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
