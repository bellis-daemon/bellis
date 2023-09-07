package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
	"log"
)

func RunTelegramBot() {
	bot, err := tgbotapi.NewBotAPI("6404196763:AAG31pmj7P6BFD4RlrIRti6DbzneJkroO1o")
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
