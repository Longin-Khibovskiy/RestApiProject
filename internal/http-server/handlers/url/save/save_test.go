package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/http-server/handlers/url/save"
	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/http-server/handlers/url/save/mocks"
	"github.com/Longin-Khibovskiy/RestApiProject.git/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{{
		name:  "Success",
		alias: "test_alias",
		url:   "https://google.com",
		// тут поля respError и mockError оставляем пустыми, тк это успешный запрос
	},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
		// другие кейсы
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewURLSaver(t)

			//	Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку,
			//	но мок должен ответить с ошибкой, к нему тоже будет запрос
			if tc.respError == "" || tc.mockError != nil {
				//	сообщаем моку, какой к нему будет запрос, и что надо вернуть
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once() // запрос будет только один
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
