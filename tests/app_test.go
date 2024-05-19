package tests

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"

	"simple_RESTapi/internal/http-server/handlers/save"
	"simple_RESTapi/internal/lib/random"
)

const (
	host = "localhost" + ":" + "6060"
)

func TestURLShortener_SavePath(t *testing.T) {
	// Универсальный способ создания URL
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	test_json := save.Request{ // Формируем тело запроса
		URL:   gofakeit.URL(),             // Генерируем случайный URL
		Alias: random.NewRandomString(10), // Генерируем случайную строку
	}

	// Создаем клиент httpexpect
	e := httpexpect.Default(t, u.String())

	e.POST("/url"). // Отправляем POST-запрос, путь - '/url'
			WithJSON(test_json).
		// WithBasicAuth("myuser", "mypass"). // Добавляем к запросу креды авторизации
		Expect().            // Далее перечисляем наши ожидания от ответа
		Status(200).         // Код должен быть 200
		JSON().Object().     // Получаем JSON-объект тела ответа
		ContainsKey("alias") // Проверяем, что в нём есть ключ 'alias'
}

// func TestURLShortener_DeletePath(t *testing.T) {
// 	// Универсальный способ создания URL
// 	u := url.URL{
// 		Scheme: "http",
// 		Host:   host,
// 	}

// 	// Создаем клиент httpexpect
// 	e := httpexpect.Default(t, u.String())

// 	e.POST("/delete"). // Отправляем POST-запрос, путь - '/url'
// 				WithJSON(delete.request{}).
// 		// WithBasicAuth("myuser", "mypass"). // Добавляем к запросу креды авторизации
// 		Expect().            // Далее перечисляем наши ожидания от ответа
// 		Status(200).         // Код должен быть 200
// 		JSON().Object().     // Получаем JSON-объект тела ответа
// 		ContainsKey("alias") // Проверяем, что в нём есть ключ 'alias'

// }
