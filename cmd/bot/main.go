/*
 *   Copyright (c) 2022 CRT_HAO 張皓鈞
 *   All rights reserved.
 *   CISH Robotics Team
 */

package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("已登入 %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // 如果收到訊息
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Document != nil {
				if update.Message.Caption != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("「%s」\n\n收到了！我會幫你記住～", update.Message.Caption))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				} else {
					var keyboard = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("跟檔名一樣", "跟檔名一樣"),
						),
					)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "這是什麼檔案？\n(回覆我檔案名稱)")
					msg.ReplyMarkup = keyboard
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "阿巴阿巴～")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
