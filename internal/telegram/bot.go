package telegram

import (
	"context"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go-inventory-management/internal/app"
)

func Start(ctx context.Context, token string, application *app.App) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf("Telegram bot authorized as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := handleUpdate(ctx, bot, application, update); err != nil {
				log.Printf("handle update: %v", err)
			}
		}
	}
}

func handleUpdate(
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	application *app.App,
	update tgbotapi.Update,
) error {
	var chatID int64
	var input app.Input
	var isStart bool

	switch {
	case update.CallbackQuery != nil:
		chatID = update.CallbackQuery.Message.Chat.ID
		input.CallbackID = update.CallbackQuery.Data
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := bot.Request(callback); err != nil {
			log.Printf("callback ack: %v", err)
		}

	case update.Message != nil:
		chatID = update.Message.Chat.ID
		text := strings.TrimSpace(update.Message.Text)
		if text == "/start" || strings.HasPrefix(text, "/start ") {
			isStart = true
		} else {
			input.Text = text
		}

	default:
		return nil
	}

	var (
		reply *app.Reply
		err   error
	)

	if isStart {
		reply, err = application.HandleStart(ctx, chatID)
	} else {
		reply, err = application.HandleInput(ctx, chatID, input)
	}
	if err != nil {
		return err
	}

	return sendReply(bot, chatID, reply)
}

func sendReply(bot *tgbotapi.BotAPI, chatID int64, reply *app.Reply) error {
	if reply == nil {
		return nil
	}

	if reply.Status != "" {
		statusMsg := tgbotapi.NewMessage(chatID, reply.Status)
		if _, err := bot.Send(statusMsg); err != nil {
			log.Printf("send status: %v", err)
		}
	}

	msg := tgbotapi.NewMessage(chatID, reply.Text)
	if len(reply.Buttons) > 0 {
		msg.ReplyMarkup = buildKeyboard(reply.Buttons)
	}

	_, err := bot.Send(msg)
	return err
}

func buildKeyboard(buttons []app.Button) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(truncate(button.Label, 64), button.ID),
		))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
