package handler

import (
	"context"
	"github/fk_reminder_bot/model"
	"github/fk_reminder_bot/ui/datetimepicker"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ReminderHandler struct {
	reminderManager *model.ReminderManager
}

func NewReminderHandler(reminderManager *model.ReminderManager) *ReminderHandler {
	return &ReminderHandler{reminderManager: reminderManager}
}

func (rh *ReminderHandler) AddReminder(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datetimepicker.New(b, AddReminderCallback)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select date",
		ReplyMarkup: kb,
	})
}

// functions callback message
func AddReminderCallback(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage, date time.Time) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You select " + date.Format("2006-01-02 15:04"),
	})
}
