package main

import (
    "log"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    // Инициализация бота
    bot, err := tgbotapi.NewBotAPI("7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM")
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true
    log.Printf("Бот запущен: %s", bot.Self.UserName)

    // Получение обновлений от Telegram
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.GetUpdatesChan(u)

    // Обработка обновлений
    for update := range updates {
        if update.Message != nil { // Проверяем, что это сообщение
            if update.Message.IsCommand() { // Если это команда
                switch update.Message.Command() {
                case "start":
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Чтобы получить доступ в приватный чат, пожалуйста, оплатите доступ.")
                    bot.Send(msg)
                case "pay":
                    // Условная проверка "оплаты"
                    if update.Message.Text == "оплачено" {
                        chatID := int64(-123456789) // Заменить на реальный ID приватного чата
                        userID := update.Message.From.ID

                        // Создание ссылки на чат
                        inviteLink, err := bot.CreateChatInviteLink(tgbotapi.CreateChatInviteLinkConfig{
                            ChatID: chatID,
                            Name:   "Доступ по оплате",
                        })
                        if err != nil {
                            log.Println("Ошибка создания ссылки:", err)
                            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка доступа в чат. Попробуйте позже.")
                            bot.Send(msg)
                        } else {
                            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Оплата подтверждена! Вот ваша ссылка в чат: "+inviteLink.InviteLink)
                            bot.Send(msg)
                        }
                    } else {
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Оплата не подтверждена. Пожалуйста, проверьте данные.")
                        bot.Send(msg)
                    }
                }
            }
        }
    }
}