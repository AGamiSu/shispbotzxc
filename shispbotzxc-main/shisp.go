package main

import (
	"log"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	logChatID    int64 = -1002496070860 // ID чата для логов
	party        = make(map[int64][]string) // Список каналов с участниками
	partyLock    sync.Mutex                 // Мьютекс для синхронизации доступа к party
	partyMessage = make(map[int64]int)      // ID последнего сообщения для каждого чата
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
        handleBattlecup(bot, msg.Chat.ID) // Добавляем вызов
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

// Обработка команды /battlecup
func handleBattlecup(bot *tgbotapi.BotAPI, chatID int64) {
    // Отправляем сообщение с кнопками для выбора ролей
    buttons := [][]tgbotapi.InlineKeyboardButton{
        {tgbotapi.NewInlineKeyboardButtonData("Керри", "role_carry")},
        {tgbotapi.NewInlineKeyboardButtonData("Мидер", "role_mid")},
        {tgbotapi.NewInlineKeyboardButtonData("Оффлейн", "role_offlane")},
        {tgbotapi.NewInlineKeyboardButtonData("Саппорт 4", "role_pos4")},
        {tgbotapi.NewInlineKeyboardButtonData("Саппорт 5", "role_pos5")},
    }

    msg := tgbotapi.NewMessage(chatID, "Собираем команду для BattleCup! Выберите роль:")
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

    // Отправляем сообщение
    sentMsg, err := bot.Send(msg)
    if err != nil {
        log.Printf("Ошибка отправки сообщения для BattleCup: %v", err)
        return
    }

    // Сохраняем ID сообщения, чтобы обновлять его
    partyMessage[chatID] = sentMsg.MessageID
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	data := callback.Data                 // Получаем данные из нажатой кнопки
	username := callback.From.UserName   // Имя пользователя
	chatID := callback.Message.Chat.ID   // ID чата

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

	// Обновление сообщения с участниками и ролями
func updateBattlecupMessage(bot *tgbotapi.BotAPI, chatID int64) {
	// Названия ролей
	roles := map[string]string{
		"role_carry":    "Керри",
		"role_mid":      "Мидер",
		"role_offlane":  "Оффлейн",
		"role_pos4":     "Саппорт 4",
		"role_pos5":     "Саппорт 5",
	}

	// Инициализируем текст для каждой роли
	roleTexts := make(map[string]string)
	for key, name := range roles {
		roleTexts[key] = name + ": ___" // Заполнено прочерками по умолчанию
	}

	// Заполняем роли участниками
	for username, role := range battlecupParticipants {
		if _, exists := roleTexts[role]; exists {
			roleTexts[role] = roles[role] + ": @" + username
		}
	}

	// Формируем итоговый текст сообщения
	messageText := "Собираем команду для BattleCup!\n\n"
	for _, role := range roles {
		messageText += roleTexts[role] + "\n"
	}

	// Добавляем кнопку для отказа от участия
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("Отказаться от участия", "leave_party")},
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Обновляем сообщение
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, partyMessage[chatID], messageText, keyboard)

	// Отправляем обновление
	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Ошибка обновления сообщения: %v", err)
	}
}

	// Обновляем сообщение для join/leave
	messageText := "Собираем пати! Текущие участники:\n" + formatPartyList(chatID)
	keyboard := generatePartyButtons(chatID)
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, partyMessage[chatID], messageText, keyboard)

	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("Ошибка обновления сообщения: %v", err)
	}
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



// Отправка логов в лог-чат
func sendLog(bot *tgbotapi.BotAPI, message string) {
	msg := tgbotapi.NewMessage(logChatID, message)
	bot.Send(msg)
}