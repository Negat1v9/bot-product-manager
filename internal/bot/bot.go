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

type Bot struct {
	timeOut int
	offset  int
	client  *tgbotapi.BotAPI
	logger  *slog.Logger
	hub     manager.Manager
}

func New(client *tgbotapi.BotAPI, timeOut int, offset int) *Bot {
	return &Bot{
		timeOut: timeOut,
		offset:  offset,
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

	b.hub = NewHub(store)

	configBot := tgbotapi.NewUpdate(b.offset)

	configBot.Timeout = b.timeOut
	// channel for getting updates
	updates := b.client.GetUpdatesChan(configBot)
	// set commands for bot
	b.client.Send(tgbotapi.NewSetMyCommands(botCommands...))

	return b.startPooling(updates)
}

// TODO: Edit parametrs in error logger, make it
func (b *Bot) startPooling(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		// Message update
		start := time.Now()
		if update.Message != nil {
			// manage cmd or message to hub
			msg, err := b.hub.MessageUpdate(update.Message)
			if err != nil {
				b.logger.Error("Manage - Create Message:", slog.String("error", err.Error()))
				continue
			}
			if _, err := b.client.Send(msg); err != nil {
				b.logger.Error("Not sending message:", slog.String("error", err.Error()))
			}
		}
		if update.CallbackQuery != nil {
			msg, err := b.hub.CallBackUpdate(*update.CallbackQuery)
			if err != nil {
				b.logger.Error("Manage callback update:", slog.String("error", err.Error()))
				continue
			}
			if _, err = b.client.Send(msg); err != nil {
				b.logger.Error("Not sending message from callback", slog.String("error", err.Error()))
			}
		}
		b.logger.Info("Query", slog.Duration("time", time.Now().Sub(start)))
	}
	return nil
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
