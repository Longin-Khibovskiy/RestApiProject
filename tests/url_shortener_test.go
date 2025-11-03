package tests

import (
	"net/url"
	"testing"

	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/http-server/handlers/url/save"
	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const host = "localhost:8082"

func TestURLShortener_HappyPath(t *testing.T) {
	//	универсальный способ создания URL
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	//	Создаем клиент httpexpect
	e := httpexpect.Default(t, u.String())

	e.Post("/url"). // Отправляем POST запрос, путь /url
		WithJSON(save.Request{ // Формируем тело запроса
			URL:   gofakeit.URL(),             // генерируем случайный URL
			Alias: random.NewRandomString(10), // Генерируем случайную строку
		}).
		WithBasicAuth("user", "pass"). // Добавляем к запросу креды авторизации
		Expect(). // Перечисляем наши ожидания от ответа
		Status(200). // Код должен быть 200
		JSON().Object(). // Получаем JSON-объект тела ответа
		ContainsKey("alias") // Проверяем, что в нем есть ключ 'alias'
}
