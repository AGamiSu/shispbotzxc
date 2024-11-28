package main

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	logChatID int64 = -1002496070860 // ID чата для логов
)

func main() {
	// Создание бота
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Printf("Бот авторизован как %s", bot.Self.UserName)

	// Обновления от Telegram
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Если пришло сообщение
			handleMessage(bot, update.Message)
		}

		if update.CallbackQuery != nil { // Обработка нажатий на кнопки
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

// Обработка текстовых сообщений
func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// Логируем входящее сообщение в лог-чат
	sendLog(bot, "Получено сообщение: "+msg.Text+" от "+msg.From.UserName)

	// Обработка команд
	switch msg.Command() {
	case "party":
		handleParty(bot, msg.Chat.ID)
	case "battlecup":
		handleBattlecup(bot, msg.Chat.ID)
	default:
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная команда"))
	}
}

// Обработка команды /party
func handleParty(bot *tgbotapi.BotAPI, chatID int64) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Вступить", "join_party")},
		{tgbotapi.NewInlineKeyboardButtonData("Выйти", "leave_party")},
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, "Собираем пати! Нажмите, чтобы присоединиться:")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// Обработка команды /battlecup
func handleBattlecup(bot *tgbotapi.BotAPI, chatID int64) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Керри", "role_carry")},
		{tgbotapi.NewInlineKeyboardButtonData("Мидер", "role_mid")},
		{tgbotapi.NewInlineKeyboardButtonData("Оффлейн", "role_offlane")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 4", "role_pos4")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 5", "role_pos5")},
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, "Собираем команду для BattleCup! Выберите роль:")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// Обработка callback кнопок
func handleCallback(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	var response string

	switch query.Data {
	case "join_party":
		response = "Вы присоединились к пати!"
	case "leave_party":
		response = "Вы вышли из пати."
	case "role_carry":
		response = "Вы выбрали роль: Керри."
	case "role_mid":
		response = "Вы выбрали роль: Мидер."
	case "role_offlane":
		response = "Вы выбрали роль: Оффлейн."
	case "role_pos4":
		response = "Вы выбрали роль: Саппорт 4."
	case "role_pos5":
		response = "Вы выбрали роль: Саппорт 5."
	}

	// Ответ на нажатие кнопки
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, response)
	bot.Send(msg)

	// Логируем нажатие кнопки
	sendLog(bot, "Нажата кнопка: "+query.Data+" от "+query.From.UserName)
}

// Отправка логов в лог-чат
func sendLog(bot *tgbotapi.BotAPI, message string) {
	msg := tgbotapi.NewMessage(logChatID, message)
	bot.Send(msg)
}