package datetimepicker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	cmdPrevMonth = iota
	cmdNextMonth
	cmdPrevYears
	cmdNextYears
	cmdCancel
	cmdBack
	cmdMonthClick
	cmdYearClick
	cmdNop

	cmdPeriodeClick
	cmdMinuteClick
	cmdHourClick
	cmdDayClick
	cmdSelectMonth
	cmdSelectYear
)

func (datePicker *DatePicker) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		datePicker.onError(err)
		return
	}
	if !ok {
		datePicker.onError(fmt.Errorf("callback answer failed"))
	}
}

func (datePicker *DatePicker) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	st := datePicker.decodeState(update.CallbackQuery.Data)

	switch st.cmd {
	case cmdYearClick:
		datePicker.year = st.param
		datePicker.showMain(ctx, b, update.CallbackQuery.Message)
	case cmdMonthClick:
		datePicker.month = time.Month(st.param)
		datePicker.showMain(ctx, b, update.CallbackQuery.Message)
	case cmdDayClick:
		datePicker.day = st.param
		datePicker.showSelectHour(ctx, b, update.CallbackQuery.Message)
	case cmdHourClick:
		datePicker.hour = st.param
		datePicker.showSelectMinute(ctx, b, update.CallbackQuery.Message)
	case cmdMinuteClick:
		datePicker.minute = st.param
		datePicker.showSelectPeriode(ctx, b, update.CallbackQuery.Message)
	case cmdPeriodeClick:
		datePicker.hour += st.param
		if datePicker.deleteOnSelect {
			_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			})
			if errDelete != nil {
				datePicker.onError(fmt.Errorf("failed to delete message onSelect: %w", errDelete))
			}
			b.UnregisterHandler(datePicker.callbackHandlerID)
		}
		datePicker.onSelect(
			ctx, b, update.CallbackQuery.Message,
			time.Date(
				datePicker.year,
				datePicker.month,
				datePicker.day,
				datePicker.hour,
				datePicker.minute,
				0, 0, time.Local))
	case cmdCancel:
		if datePicker.deleteOnCancel {
			_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			})
			if errDelete != nil {
				datePicker.onError(fmt.Errorf("failed to delete message onCancel: %w", errDelete))
			}
			b.UnregisterHandler(datePicker.callbackHandlerID)
		}
		datePicker.onCancel(ctx, b, update.CallbackQuery.Message)
	case cmdPrevYears:
		datePicker.showSelectYear(ctx, b, update.CallbackQuery.Message, st.param)
	case cmdNextYears:
		datePicker.showSelectYear(ctx, b, update.CallbackQuery.Message, st.param)
	case cmdPrevMonth:
		datePicker.month--
		if datePicker.month == 0 {
			datePicker.month = 12
			datePicker.year--
		}
		datePicker.showMain(ctx, b, update.CallbackQuery.Message)
	case cmdNextMonth:
		datePicker.month++
		if datePicker.month > 12 {
			datePicker.month = 1
			datePicker.year++
		}
		datePicker.showMain(ctx, b, update.CallbackQuery.Message)
	case cmdBack:
		datePicker.showMain(ctx, b, update.CallbackQuery.Message)
	case cmdSelectMonth:
		datePicker.showSelectMonth(ctx, b, update.CallbackQuery.Message)
	case cmdSelectYear:
		datePicker.showSelectYear(ctx, b, update.CallbackQuery.Message, datePicker.year)
	case cmdNop:
		// do nothing
	default:
		datePicker.onError(fmt.Errorf("unknown command: %d", st.cmd))
	}

	datePicker.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (datePicker *DatePicker) showSelectPeriode(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage) {
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    mes.Chat.ID,
		MessageID: mes.MessageID,
		Text:      "Select Periode",
	})

	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectPeriode, %w", err))

	}
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildPeriodeKeyboard()},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectTime, %w", err))
	}
}
func (datePicker *DatePicker) showSelectMinute(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage) {
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    mes.Chat.ID,
		MessageID: mes.MessageID,
		Text:      "Select Minute",
	})

	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectMinute, %w", err))

	}
	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildMinuteKeyboard()},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectTime, %w", err))
	}
}
func (datePicker *DatePicker) showSelectHour(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage) {
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    mes.Chat.ID,
		MessageID: mes.MessageID,
		Text:      "Select Hour",
	})

	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectHour, %w", err))

	}

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildHourKeyboard()},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectTime, %w", err))
	}
}
func (datePicker *DatePicker) showSelectMonth(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage) {
	_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildMonthKeyboard()},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectMonth, %w", err))
	}
}

func (datePicker *DatePicker) showSelectYear(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage, currentYear int) {
	_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildYearKeyboard(currentYear)},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showSelectYear, %w", err))
	}
}

func (datePicker *DatePicker) showMain(ctx context.Context, b *bot.Bot, mes models.InaccessibleMessage) {
	_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      mes.Chat.ID,
		MessageID:   mes.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: datePicker.buildKeyboard()},
	})
	if err != nil {
		datePicker.onError(fmt.Errorf("error edit message in showMain, %w", err))
	}
}
