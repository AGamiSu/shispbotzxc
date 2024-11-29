package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	logChatID             int64 = -1002496070860 // ID чата для логов
	party                 = make(map[int64][]string) // Список каналов с участниками
	partyLock             sync.Mutex                 // Мьютекс для синхронизации доступа к party
	partyMessage          = make(map[int64]int)      // ID последнего сообщения для каждого чата
	battlecupParticipants = make(map[string]string)  // Словарь для хранения ролей участников
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

	// Запускаем горутину для отслеживания времени
	go scheduleReminder(bot)

	// Основной цикл для обработки обновлений
	for update := range updates {
		if update.Message != nil { // Если пришло сообщение
			handleMessage(bot, update.Message)
		}

		if update.CallbackQuery != nil { // Обработка нажатий на кнопки
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

// Отправка напоминания для подтверждения записи на BattleCup каждую субботу в 20:30
func scheduleReminder(bot *tgbotapi.BotAPI) {
	for {
		// Рассчитываем время до следующего 20:30 в субботу
		now := time.Now()
		todaySaturday := now.Add(time.Hour * time.Duration((6-int(now.Weekday()))+7)%7)
		reminderTime := time.Date(todaySaturday.Year(), todaySaturday.Month(), todaySaturday.Day(), 20, 30, 0, 0, todaySaturday.Location())

		// Если текущее время позже 20:30, настраиваем на следующую неделю
		if now.After(reminderTime) {
			reminderTime = reminderTime.Add(7 * 24 * time.Hour)
		}

		// Засыпаем до времени напоминания
		time.Sleep(reminderTime.Sub(now))

		// Напоминаем всем участникам
		sendBattlecupReminder(bot)
	}
}

// Отправка напоминания о подтверждении записи на BattleCup
func sendBattlecupReminder(bot *tgbotapi.BotAPI) {
	// Формируем текст сообщения с участниками и их ролями
	messageText := "Время подтвердить участие в BattleCup! Текущие участники:\n"

	// Проверяем участников и добавляем их в сообщение
	for username, role := range battlecupParticipants {
		messageText += role + " " + username + "\n"
	}

	// Отправляем сообщение в чат
	for chatID := range party {
		msg := tgbotapi.NewMessage(chatID, messageText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}
	}
}

// Обработка текстовых сообщений
func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// Проверяем, является ли сообщение командой
	if !msg.IsCommand() {
		return // Если сообщение не команда, игнорируем
	}

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

	// Отправка нового сообщения
	messageText := "Собираем пати! Текущие участники:\n" + formatPartyList(chatID)
	keyboard := generatePartyButtons(chatID)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard

	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
		return
	}

	partyMessage[chatID] = sentMsg.MessageID // Сохраняем ID нового сообщения
}

// Форматирование списка участников пати
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

// Обработка callback кнопок
func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	username := callback.From.UserName
	chatID := callback.Message.Chat.ID

	partyLock.Lock()
	defer partyLock.Unlock()

	if strings.HasPrefix(data, "join_") {
		if !isInParty(chatID, username) {
			party[chatID] = append(party[chatID], username)
		}
	} else if strings.HasPrefix(data, "leave_") {
		party[chatID] = removeUserFromParty(party[chatID], username)
	} else if strings.HasPrefix(data, "role_") {
		// Обработка выбора роли в BattleCup
		battlecupParticipants[username] = data
	}

	// Обновляем текст и кнопки
	messageText := "Собираем пати! Текущие участники:\n" + formatPartyList(chatID)
	keyboard := generatePartyButtons(chatID)
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, partyMessage[chatID], messageText, keyboard)

	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Ошибка обновления сообщения: %v", err)
	}

	// Ответ на CallbackQuery через Request
	callbackResponse := tgbotapi.NewCallback(callback.ID, "Действие выполнено")
	if _, err := bot.Request(callbackResponse); err != nil {
		log.Printf("Ошибка ответа на CallbackQuery: %v", err)
	}
}

// Проверка, находится ли пользователь в пати
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

// Обработка команды /battlecup
func handleBattlecup(bot *tgbotapi.BotAPI, chatID int64) {
	// Кнопки для ролей
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Керри", "role_carry")},
		{tgbotapi.NewInlineKeyboardButtonData("Мидер", "role_mid")},
		{tgbotapi.NewInlineKeyboardButtonData("Оффлейн", "role_offlane")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 4", "role_pos4")},
		{tgbotapi.NewInlineKeyboardButtonData("Саппорт 5", "role_pos5")},
	}

	// Формируем и отправляем сообщение
	messageText := "Выберите свою роль для BattleCup:"
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	bot.Send(msg)
}

// Функция логирования сообщений
func sendLog(bot *tgbotapi.BotAPI, message string) {
	logMessage := tgbotapi.NewMessage(logChatID, message)
	bot.Send(logMessage)
}