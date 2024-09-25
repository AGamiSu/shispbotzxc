package main

import (
	"log"
	"os"

	tgbotapi "7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM"
)

func main() {
	// Получаем токен из переменных окружения (или можно напрямую вставить)
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Panic("TELEGRAM_BOT_TOKEN environment variable is not set!")
	}

	// Инициализируем бота с помощью токена
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Включаем отладочный режим (если нужно)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Создаем новый апдейт канал
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Запускаем получение апдейтов
	updates := bot.GetUpdatesChan(u)

	// Обрабатываем каждый апдейт
	for update := range updates {
		if update.Message != nil { // Игнорируем любые обновления без сообщения
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Отправляем обратно сообщение с приветствием
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я твой бот для приватных чатов.")
			bot.Send(msg)
		}
	}
}
