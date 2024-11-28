package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	partyMessageID int
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "пати":
				handleParty(bot, update)
			case "батлкап":
				handleBattleCup(bot)
			}
		}
	}
}

func handleParty(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Собираем игроков! Нажмите кнопку, чтобы присоединиться.")
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Присоединиться", "join_party"),
			tgbotapi.NewInlineKeyboardButtonData("Покинуть", "leave_party"),
		},
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ReplyMarkup = keyboard

	// Удаление предыдущего сообщения
	if partyMessageID != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, partyMessageID)
		bot.Send(deleteMsg)
	}

	sentMsg, _ := bot.Send(msg)
	partyMessageID = sentMsg.MessageID // Сохраняем ID нового сообщения
}


func handleBattleCup(bot *tgbotapi.BotAPI) {
	loc, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(loc)
	nextSaturday := now.AddDate(0, 0, (6-int(now.Weekday()))%7)
	scheduleTime := time.Date(nextSaturday.Year(), nextSaturday.Month(), nextSaturday.Day(), 9, 0, 0, 0, loc)

	// Таймер на утреннее сообщение
	go func() {
		time.Sleep(time.Until(scheduleTime.Add(-12 * time.Hour))) // 9 утра МСК
		sendBattleCupMessage(bot)
	}()

	// Таймер на пинг за 30 минут
	go func() {
		time.Sleep(time.Until(scheduleTime.Add(-30 * time.Minute)))
		sendReminder(bot)
	}()
}

func sendBattleCupMessage(bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(-1001740769275, "Выберите вашу роль для Батлкапа:")
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Керри", "carry"),
			tgbotapi.NewInlineKeyboardButtonData("Мидер", "mid"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Офлейн", "offlane"),
			tgbotapi.NewInlineKeyboardButtonData("Саппорт 4", "pos4"),
			tgbotapi.NewInlineKeyboardButtonData("Саппорт 5", "pos5"),
		},
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	bot.Send(msg)
}

func sendReminder(bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(-1001740769275, "Напоминание! Подтвердите участие в Батлкапе!")
	bot.Send(msg)
}
