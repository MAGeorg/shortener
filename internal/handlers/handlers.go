package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/core"
	customerr "github.com/MAGeorg/shortener.git/internal/errors"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/MAGeorg/shortener.git/internal/utils"
)

// структура содержащая необходимые данные для обработки запросов
// функция обработки запросов - методы структуры.
type AppHandler struct {
	a *appdata.AppData
}

// обработка POST запроса.
func (h *AppHandler) CreateHashURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	defer r.Body.Close()
	urlStr, err := io.ReadAll(r.Body)

	// проверка входящего URL.
	if err != nil || !utils.CheckURL(string(urlStr)) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ctx := context.Background()
	urlHash, err := core.CreateShotURL(
		ctx,
		&core.InputValueForWriteFile{
			Stor:        h.a.StorageURL,
			BaseAddress: h.a.BaseAddress,
			URL:         string(urlStr),
		})

	switch {
	// проверка на ошибку unique_violation.
	case errors.Is(err, customerr.ErrAccessDenied):
		w.WriteHeader(http.StatusConflict)

	// ошибка при генерации сокращенного URL, возращаем 500.
	case err != nil:
		h.a.Logger.Errorln("error create new record", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return

	// формирование положительного ответа.
	default:
		w.WriteHeader(http.StatusCreated)
	}

	_, err = w.Write([]byte(urlHash))
	if err != nil {
		// ошибка при записи ответа в Body, возращаем 500.
		h.a.Logger.Errorln("error when write answer", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// обработка GET запросв.
func (h *AppHandler) GetOriginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if len(r.URL.String()) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := context.Background()
	url, err := core.GetOriginURL(ctx, h.a.StorageURL, r.URL.String()[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа.
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// обработка POST запроса в формате JSON.
func (h *AppHandler) CreateHashURLJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// преобразуем bytes (JSON) в map.
	var urlJSON models.OriginURL
	if err := json.Unmarshal(data, &urlJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !utils.CheckURL(urlJSON.URL) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ctx := context.Background()
	urlHash, err := core.CreateShotURL(
		ctx,
		&core.InputValueForWriteFile{
			Stor:        h.a.StorageURL,
			BaseAddress: h.a.BaseAddress,
			URL:         urlJSON.URL,
		})

	switch {
	// проверка на ошибку unique_violation.
	case errors.Is(err, customerr.ErrAccessDenied):
		w.WriteHeader(http.StatusConflict)

	// ошибка при генерации сокращенного URL, возращаем 500.
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return

	// формирование положительного ответа.
	default:
		w.WriteHeader(http.StatusCreated)
	}

	resp, err := json.Marshal(models.ResponseHashURL{URL: urlHash})
	if err != nil {
		// ошибка при сериализации JSON объекта.
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		// ошибка при записи ответа в Body, возращаем 500.
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// обработка GET запроса для ping DataBase.
func (h *AppHandler) PingDB(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := core.PingDB(h.a.DSNdatabase)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// обработка POST запроса для создания сокращенных url для списка url.
func (h *AppHandler) CreateHashURLBatchJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// проверка, что на вход пришел не пустой body.
	if len(data) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// преобразуем bytes (JSON) в map.
	var batchJSON []models.DataBatch
	if err := json.Unmarshal(data, &batchJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	res, err := core.CreateShotURLBatch(ctx, h.a.StorageURL, h.a.BaseAddress, batchJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// формируем body ответа.
	w.WriteHeader(http.StatusCreated)
	var b []byte

	if len(res) > 0 {
		b, err = json.Marshal(res)
		if err != nil {
			// ошибка при сериализации JSON объекта.
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if _, err := w.Write(b); err != nil {
		// ошибка при записи ответа в Body, возращаем 500.
		w.WriteHeader(http.StatusInternalServerError)
	}
}
