package handler

import (
	"context"
	"github/fk_reminder_bot/model"
	"github/fk_reminder_bot/ui/datetimepicker"
	"regexp"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ReminderHandler struct {
	reminderManager *model.ReminderManager
	userManager     *model.UserManager
}

func NewReminderHandler(reminderManager *model.ReminderManager, userManager *model.UserManager) *ReminderHandler {
	return &ReminderHandler{reminderManager: reminderManager, userManager: userManager}
}

func (rh *ReminderHandler) AddReminder(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datetimepicker.New(b, AddReminderCallback)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select date",
		ReplyMarkup: kb,
	})
}

func (rh *ReminderHandler) SetReminder(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	regex := `^\/set\s(\d{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])\s([01][0-9]|2[0-3]):([0-5][0-9])$`

	// Compile ekspresi reguler
	re := regexp.MustCompile(regex)

	// Cek apakah string cocok dengan ekspresi reguler
	if !re.MatchString(text) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "format yang anda masukkan salah. format harus YYYY-MM-DD HH-MM",
		})
		return
	}

	userId := update.Message.From.ID
	user, _ := rh.userManager.GetByUserID(userId)
	if user.ID <= 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "silakan subscribe bot ini dengan command /start",
		})
		return
	}

	// rh.reminderManager.AddReminder(int(userId), )
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Success add reminder for date: ",
	})

}

// functions callback message
func AddReminderCallback(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage, date time.Time) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You select " + date.Format("2006-01-02 15:04"),
	})
}
