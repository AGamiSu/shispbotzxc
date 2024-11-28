package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Ваш токен
	bot, err := tgbotapi.NewBotAPI("7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM")
	if err != nil {
		log.Fatalf("Ошибка подключения к боту: %v", err)
	}

	bot.Debug = true // Включить отладочный режим для детального вывода

	// Ваш User ID
	yourTelegramID := int64(-1001740769275)

	// Пробуем отправить тестовое сообщение с кнопками
	msg := tgbotapi.NewMessage(yourTelegramID, "Это тестовое сообщение с кнопками.")
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Кнопка 1", "data1"),
			tgbotapi.NewInlineKeyboardButtonData("Кнопка 2", "data2"),
		},
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Отправляем сообщение
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение: %v", err)
	}

	log.Println("Сообщение отправлено успешно!")
}
