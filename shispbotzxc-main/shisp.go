package main

import (
    "log"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const botToken = "7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM"

func main() {
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Бот запущен: %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil { // Проверяем, что это сообщение
            if update.Message.IsCommand() { // Если это команда
                switch update.Message.Command() {
                case "start":
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я SHISPbot, чтобы получить доступ в приватный чат, пожалуйста, оплатите доступ.")
                    bot.Send(msg)
                }
            }
        }
    }
}