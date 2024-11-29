package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка команды BattleCup
func handleBattlecup(bot *tgbotapi.BotAPI, chatID int64) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Керри", "role_carry")},
		{tgbotapi.NewInlineKeyboardButtonData("Мидер", "role_mid")},
		{tgbotapi.NewInlineKeyboardButtonData("Оффлейн", "role_offlane")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 4", "role_pos4")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 5", "role_pos5")},
	}

	msg := tgbotapi.NewMessage(chatID, "Собираем команду для BattleCup! Выберите роль:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения для BattleCup: %v", err)
	}
}

// Выбор роли для BattleCup
func selectRole(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, chatID int64, username string, role string) {
	// Проверяем, выбрал ли пользователь роль
	if _, exists := battlecupParticipants[username]; exists {
		// Сообщаем пользователю, что он уже выбрал роль
		bot.Request(tgbotapi.NewCallback(callback.ID, "Вы уже выбрали роль"))
		return
	}

	// Назначаем роль пользователю
	battlecupParticipants[username] = role

	// Сообщаем о выборе
	bot.Request(tgbotapi.NewCallback(callback.ID, "Вы выбрали роль: "+role))

	// Обновляем сообщение с ролями
	updateBattlecupMessage(bot, chatID)

	battlecupParticipants[username] = role
	bot.Request(tgbotapi.NewCallback(callback.ID, "Вы выбрали роль: "+role))

	// Обновляем сообщение с ролями
	updateBattlecupMessage(bot, chatID)
}

// Обновление сообщения с ролями
func updateBattlecupMessage(bot *tgbotapi.BotAPI, chatID int64) {
	// Логика обновления сообщения с ролями
	// (реализуй аналогично тому, что ты писал в изначальном коде)
}
func LeaveBattlecup(username string) {
	delete(battlecupParticipants, username)
}
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	data := callback.Data // Извлечение данных кнопки
	username := callback.From.UserName
	chatID := callback.Message.Chat.ID
	
	partyLock.Lock()
	defer partyLock.Unlock()

	if strings.HasPrefix(data, "join_") {
		// Добавление пользователя в список
		if !isInParty(chatID, username) {
			party[chatID] = append(party[chatID], username)
		}
	} else if strings.HasPrefix(data, "leave_") {
		// Удаление пользователя из списка
		party[chatID] = removeUserFromParty(party[chatID], username)
	} else if strings.HasPrefix(data, "role_") {
		// Выбор роли для BattleCup
		if _, exists := battlecupParticipants[username]; exists {
			// Сообщаем, что роль уже выбрана
			bot.Request(tgbotapi.NewCallback(callback.ID, "Вы уже выбрали роль"))
			return
		}

		// Добавляем выбранную роль
		battlecupParticipants[username] = data
		bot.Request(tgbotapi.NewCallback(callback.ID, "Вы выбрали роль: "+data))

		// Обновляем сообщение с ролями
		updateBattlecupMessage(bot, chatID)
	} else if data == "leave_party" {
		// Удаление роли пользователя
		delete(battlecupParticipants, username)
		bot.Request(tgbotapi.NewCallback(callback.ID, "Вы отказались от участия"))
		updateBattlecupMessage(bot, chatID)
	}
}