package bot

import (
	"log"

	"github.com/bellis-daemon/bellis/common/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
)

func RunTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(storage.Firebase().ConfigGetString("telegram_bot_token"))
	if err != nil {
		glgf.Error(err)
		return
	}
	glgf.Info("Telegram authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60
	bot.GetUpdatesChan(u)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
