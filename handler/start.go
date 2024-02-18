package handler

import (
	"context"
	"github/fk_reminder_bot/model"
	"github/fk_reminder_bot/ui/datetimepicker"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StartHandler struct {
	userManager *model.UserManager
}

func NewStartHandler(userManager *model.UserManager) *StartHandler {
	return &StartHandler{userManager: userManager}
}

func (sh *StartHandler) StartBot(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datetimepicker.New(b, StartBotCallback)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select date",
		ReplyMarkup: kb,
	})
}

// functions callback message
func StartBotCallback(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage, date time.Time) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You select " + date.Format("2006-01-02 15:04"),
	})
}
