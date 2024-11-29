package main

import (
	"strconv"
	"log"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func handleParty(bot *tgbotapi.BotAPI, chatID int64) {
	partyLock.Lock()
	defer partyLock.Unlock()

	if _, exists := party[chatID]; !exists {
		party[chatID] = []string{} // Инициализируем канал, если его ещё нет
	}

	// Удаление старого сообщения
	if msgID, exists := partyMessage[chatID]; exists {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		bot.Request(deleteMsg)
	}
}

// Вход в пати
func joinParty(bot *tgbotapi.BotAPI, chatID int64, username string) {
	if !isInParty(chatID, username) {
		party[chatID] = append(party[chatID], username)
		updatePartyMessage(bot, chatID)
	}
}

// Выход из пати
func leaveParty(bot *tgbotapi.BotAPI, chatID int64, username string) {
	party[chatID] = removeUserFromParty(party[chatID], username)
	updatePartyMessage(bot, chatID)
}

// Проверка, что пользователь уже в пати
func isInParty(chatID int64, username string) bool {
	for _, user := range party[chatID] {
		if user == username {
			return true
		}
	}
	return false
}

// Удаление пользователя из пати
func removeUserFromParty(users []string, username string) []string {
	for i, user := range users {
		if user == username {
			return append(users[:i], users[i+1:]...)
		}
	}
	return users
}

// Обновление сообщения с участниками
func updatePartyMessage(bot *tgbotapi.BotAPI, chatID int64) {
	messageText := "Собираем пати! Текущие участники:\n" + formatPartyList(chatID)
	keyboard := generatePartyButtons(chatID)
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, partyMessage[chatID], messageText, keyboard)

	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Ошибка обновления сообщения: %v", err)
	}
}

// Форматирование списка участников
func formatPartyList(chatID int64) string {
	participants := party[chatID]
	if len(participants) == 0 {
		return "Никто ещё не присоединился."
	}
	return "- " + strings.Join(participants, "\n- ")
}

// Генерация кнопок для пати
func generatePartyButtons(chatID int64) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Войти", "join_"+strconv.FormatInt(chatID, 10))},
		{tgbotapi.NewInlineKeyboardButtonData("Выйти", "leave_"+strconv.FormatInt(chatID, 10))},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}