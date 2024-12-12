package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/telebot.v3"
)

const (
	apiURL   = "https://api.openweathermap.org/data/2.5/weather?q=Moscow&appid=6f2d0398e8f36240d62f83eda6337f1a&units=metric"
	botToken = "7546019211:AAF5pJWX6vmusn-7u1Bp_NjMvAiYaKqmNQw"
)

// WeatherResponse структура для ответа OpenWeatherMap
type WeatherResponse struct {
	Sys struct {
		Sunrise int64 `json:"sunrise"` // UNIX время рассвета
	} `json:"sys"`
}

// getSunriseTime получает время рассвета из OpenWeatherMap
func getSunriseTime() (time.Time, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка при запросе API: %v", err)
	}
	defer resp.Body.Close()

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return time.Time{}, fmt.Errorf("ошибка при декодировании ответа API: %v", err)
	}

	// Конвертируем время рассвета из UNIX в time.Time
	return time.Unix(data.Sys.Sunrise, 0), nil
}

// scheduleMessage планирует отправку сообщения
func scheduleMessage(bot *telebot.Bot, chatID int64, sunriseTime time.Time) {
	// Рассчитываем, сколько времени ждать до рассвета
	duration := time.Until(sunriseTime)

	if duration <= 0 {
		log.Println("Рассвет уже прошел, планируем на следующий день.")
		return
	}

	log.Printf("Сообщение будет отправлено в %v\n", sunriseTime)

	time.Sleep(duration)

	msg := fmt.Sprintf("Время рассвета: %s", sunriseTime.Format("15:04"))
	bot.Send(&telebot.User{ID: chatID}, msg)
}

func main() {
	// Настройки бота
	pref := telebot.Settings{
		Token:  botToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	// ID чата — заменить на нужный
	chatID := int64(12345678)

	// Команда /start для подписки
	bot.Handle("/start", func(c telebot.Context) error {
		go func() {
			for {
				// Получаем время рассвета
				sunriseTime, err := getSunriseTime()
				if err != nil {
					log.Printf("Ошибка получения времени рассвета: %v", err)
					return
				}

				// Планируем сообщение
				scheduleMessage(bot, chatID, sunriseTime)

				// Ждем следующего дня
				time.Sleep(24 * time.Hour)
			}
		}()
		return c.Send("Вы подписались на уведомления о времени рассвета!")
	})

	// Запуск бота
	log.Println("Бот запущен!")
	bot.Start()
}