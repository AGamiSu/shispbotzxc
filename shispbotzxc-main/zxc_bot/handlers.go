package main

import (
	
	//"strings"
	//"strconv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработка сообщений
func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	switch msg.Command() {
	case "party":
		handleParty(bot, msg.Chat.ID)
	case "battlecup":
		handleBattlecup(bot, msg.Chat.ID)
	default:
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная команда"))
	}
}