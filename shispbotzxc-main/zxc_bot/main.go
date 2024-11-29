package main

import (
	"log"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	battlecupParticipants = make(map[string]string)  // Хранение ролей для участников BattleCup
	partyMessage          = make(map[int64]int)      // ID сообщений для чатов
	party                 = make(map[int64][]string) // Участники пати
	partyLock sync.Mutex
)


func main() {
	// Создание бота
	bot, err := tgbotapi.NewBotAPI("7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM")
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}
	log.Printf("Бот авторизован как %s", bot.Self.UserName)
	

	// Обновления от Telegram
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		}
		if update.CallbackQuery != nil {
			HandleCallback(bot, update.CallbackQuery)
		}
	}
}