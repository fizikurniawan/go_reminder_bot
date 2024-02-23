package main

import (
	"context"
	"github/fk_reminder_bot/config"
	"github/fk_reminder_bot/handler"
	"github/fk_reminder_bot/model"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
)

// Send any text message to the bot after the bot has been started

func main() {
	botToken := config.GoDotEnvVariable("API_TOKEN")
	dbName := config.GoDotEnvVariable("DB_NAME")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize bot options
	opts := []bot.Option{}
	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	db, err := config.InitDB(dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize the reminder manager
	userManager := model.NewUserManager(db)
	reminderManager := model.NewReminderManager(db)

	// Initialize the reminder handler
	reminderHandler := handler.NewReminderHandler(reminderManager, userManager)
	startHandler := handler.NewStartHandler(userManager)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler.StartBot)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypeExact, reminderHandler.AddReminder)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/set", bot.MatchTypeContains, reminderHandler.SetReminder)

	// Start bot
	go func() {
		b.Start(ctx)
	}()

	// Start scheduler to check reminders periodically
	go func() {
		for {
			time.Sleep(1 * time.Second) // Adjust the interval as needed
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
}
