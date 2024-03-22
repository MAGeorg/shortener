package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonHandle(t *testing.T) {
	// структура для хранения ожидаемых значений
	type want struct {
		code        int
		body        string
		contentType string
	}

	// структура для хранения данных для запросов
	type request struct {
		method      string
		contentType string
		body        string
		hashURL     string
	}

	// создаем набор тестовых данных для запроса и проверки ответа
	tests := []struct {
		name string
		req  request
		want want
	}{
		{
			name: "negative GET test 1",
			req: request{
				method:      http.MethodGet,
				hashURL:     "",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusNotFound,
				body:        "",
				contentType: "text/plain",
			},
		},
		{
			name: "negative GET test 2",
			req: request{
				method:      http.MethodGet,
				hashURL:     "12345",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusNotFound,
				body:        "",
				contentType: "text/plain",
			},
		},
		{
			name: "positive POST test 1",
			req: request{
				method:      http.MethodPost,
				body:        "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusCreated,
				body:        "/759827921",
				contentType: "text/plain",
			},
		},
		{
			name: "repeat positive POST test 1",
			req: request{
				method:      http.MethodPost,
				body:        "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusCreated,
				body:        "/759827921",
				contentType: "text/plain",
			},
		},
		{
			name: "positive GET test 1",
			req: request{
				method:      http.MethodGet,
				hashURL:     "/759827921",
				contentType: "text/plain",
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				body:        "https://practicum.yandex.ru/",
				contentType: "text/plain",
			},
		},
	}

	// так как будем использовать много assert заведем свой Assertions object
	asserts := assert.New(t)

	// запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var req *http.Request

			if test.req.method == http.MethodGet {
				req = httptest.NewRequest(test.req.method, fmt.Sprintf("/%s", test.req.hashURL), strings.NewReader(""))
			} else {
				req = httptest.NewRequest(test.req.method, "/", strings.NewReader(test.req.body))
			}

			// заполняем необходимые поля и выставляем ResponseRecorder для записи ответа сервера
			req.Header.Set("Content-Type", test.req.contentType)
			w := httptest.NewRecorder()
			CommonHandle(w, req)

			result := w.Result()
			defer result.Body.Close()

			asserts.Equal(test.want.code, result.StatusCode)
			asserts.Equal(test.want.contentType, result.Header.Get("Content-Type"))

			switch test.req.method {
			case http.MethodGet:
				asserts.Contains(result.Header.Get("location"), test.want.body)
			case http.MethodPost:
				asserts.Contains(w.Body.String(), test.want.body)
			}

		})
	}
}
