package handler

import (
	"context"
	"fmt"
	"github/fk_reminder_bot/model"
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
	from := update.Message.From
	user, err := sh.userManager.GetByUserID(from.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if user.UserID > 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You has already subscribe bot. type /help to show all features",
		})
		return
	}

	var u model.User
	u.FirstName = from.FirstName
	u.LastName = from.LastName
	u.IsActive = true
	u.UserID = int(from.ID)
	u.JoinAt = time.Now()

	_, err = sh.userManager.AddUser(u)
	if err != nil {
		fmt.Println("error add user")
		fmt.Println(err.Error())
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to FK Reminder Bot. type '/help' for show all features",
	})
}
