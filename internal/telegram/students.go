package telegram

import (
	"log"
	"log/slog"

	"github.com/kostrominoff/go-pgtk-schedule/internal/schedule"
	"github.com/kostrominoff/go-pgtk-schedule/internal/storage"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type StudentsBot struct {
	Storage  *storage.Database
	Schedule *schedule.Schedule
	Bot      *telego.Bot
}

func NewStudentsBot(
	token string,
	storage *storage.Database,
	schedule *schedule.Schedule,
) *StudentsBot {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	return &StudentsBot{
		Storage:  storage,
		Schedule: schedule,
		Bot:      bot,
	}
}

func (b *StudentsBot) Start() {
	updates, err := b.Bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatal(err)
	}

	bh, _ := th.NewBotHandler(b.Bot, updates)

	defer bh.Stop()
	defer b.Bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, err := bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), "Hello, world"))
		if err != nil {
			slog.Error(err.Error())
		}
	}, th.CommandEqual("start"))

	bh.Start()
}
