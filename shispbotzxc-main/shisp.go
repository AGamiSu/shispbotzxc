package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Константы
const requiredAmount = 123    // Минимальная сумма доната
const chatID = -1001740769275 // ID приватной группы
const botToken = "7625350088:AAGY2qai8kYZ9cwbowBODmOtFlwjhoO8ubM"
const donationToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxMzcyNiIsImp0aSI6ImQyYmE5YTExZmE0Mjc4MDVmNTNhYzNlN2VjODgzNTcyZDhhNjg2OGY4N2NiN2ExZjU1MjhkMmIzYzExZWU3MjYxYjdkMGQ0Zjk0OTdlOGI0IiwiaWF0IjoxNzMyMzEzNjQzLjE2NjcsIm5iZiI6MTczMjMxMzY0My4xNjY3LCJleHAiOjE3NjM4NDk2NDMuMTU3Mywic3ViIjoiIiwic2NvcGVzIjpbXX0.VLrW-7pNHQrEu-Hvbpc-kW7ThV_r6vGrhEno3HyO-Uh6EkThtKtzVlgc8BryUClYdF3ZWluukO8nlA5H6HBOax0A7uHaWIRa0CDX3HRDSGdIw_f5V-QTw4FzeiU5szM-1A23xSqnbuTns_ZJlhCQK988Mwo_IjHO0LlIa1BapFhzZxdH6WpYxFwT-8UhVqthOLNCe-5P0JagUhqP_nfHgSQorFuButSDzgKZ115h6P5KHD2OgC0MspvNxIMxW11Z4ndXHux2_GH1hJHKXjawZbcwnF02qUYDWgMh3Lmx0c5_uNuF4ps4HBPcRCrAOlVTpgq8BOL9jg6mT2566aPje3n6fY_BbdTdJjI4CsoITnJuU8B4Bz_XAnWLOXgGsTc0axxgo3dGnXwPbLfBUr7IYd1_FFJmxA7iWrvstyHjr1sbd6ZFC6LVNdzROkOYjWqf7exY6V8xZZwDzSJdWoqtkRCVoJ1EaaUrTusatKXMrVmkmzCQE97rbM7Y7jNLOEU6qrVWuDCIfR47z8fi67ter13gEmt1Nulm5KzwkAHCqtWPDZ9L6OVPOwmQXwl3FQNcc8os2dBeku48WlvvIlOtO7NrBKnQSD2JsifWSI6L4kPC8QLN35KY5Du62hLDJ73mpD0OvEwE4tepycDSvdzcxHn4pY6t8s69iroYqpgAzHg" // Ваш токен DonationAlerts

// Структура для данных доната
type Donation struct {
	Username string  `json:"username"`
	Message  string  `json:"message"`
	Amount   float64 `json:"amount"`
}

type DonationsResponse struct {
	Data []Donation `json:"data"`
}

func main() {
	// Инициализация Telegram бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я SHISPbot. Чтобы получить доступ в приватный чат, отправьте донат на сумму 123 рубля и укажите ваш @тег Telegram в сообщении.")
				bot.Send(msg)

			case "pay":
				// Получаем список донатов и проверяем
				userTag := update.Message.From.UserName
				if userTag == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас не установлен @тег Telegram. Установите его в настройках профиля и повторите команду /pay.")
					bot.Send(msg)
					continue
				}

				donationAmount, err := checkDonation(userTag, donationToken)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка проверки доната: "+err.Error())
					bot.Send(msg)
					continue
				}

				if donationAmount >= requiredAmount {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за донат! Вы получили доступ в группу Приватка Шиспа. Вот ссылка на группу: https://t.me/privatka_shispa")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Оплата не подтверждена или сумма меньше 123 рублей. Пожалуйста, проверьте данные.")
					bot.Send(msg)
				}
			}
		}
	}
}

// Функция для проверки доната
func checkDonation(userTag, token string) (float64, error) {
	url := "https://www.donationalerts.com/api/v1/alerts/donations"

	// Создаем HTTP-запрос с токеном
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+donationToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ошибка API DonationAlerts: %s", resp.Status)
	}

	// Декодируем ответ
	var donationsResp DonationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&donationsResp); err != nil {
		return 0, err
	}

	// Ищем подходящий донат
	for _, donation := range donationsResp.Data {
		if strings.Contains(donation.Message, "@"+userTag) && donation.Amount >= requiredAmount {
			return donation.Amount, nil
		}
	}

	return 0, fmt.Errorf("подходящий донат не найден")
}
