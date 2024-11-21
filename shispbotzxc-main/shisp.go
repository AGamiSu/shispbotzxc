package main

import (
    "log"
    "os"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/joho/godotenv"
)

func main() {
    // Загружаем переменные окружения
    err := godotenv.Load("/Users/agamisu/golang-codes/shispbotzxc-main/.env.save")
    if err != nil {
        log.Fatal("Ошибка загрузки .env файла")
    }

    botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Бот запущен: %s", bot.Self.UserName)
}
