package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/core"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/MAGeorg/shortener.git/internal/utils"
)

// структура содержащая необходимые данные для обработки запросов
// функция обработки запросов - методы структуры
type AppHandler struct {
	a *appdata.AppData
}

// обработка POST запроса
func (h *AppHandler) CreateHashURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	defer r.Body.Close()
	urlStr, err := io.ReadAll(r.Body)

	// проверка входящего URL
	if err != nil || !utils.CheckURL(string(urlStr)) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	urlHash, err := core.CreateShotURL(&core.InputValueForWriteFile{
		Stor:        h.a.StorageURL,
		BaseAddress: h.a.BaseAddress,
		URL:         string(urlStr),
	})

	if err != nil {
		// ошибка при генерации сокращенного URL, возращаем 500
		h.a.Logger.Errorln("error create new record", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(urlHash))
	if err != nil {
		// ошибка при записи ответа в Body, возращаем 500
		h.a.Logger.Errorln("error when write answer", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// обработка GET запросв
func (h *AppHandler) GetOriginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if len(r.URL.String()) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := core.GetOriginURL(h.a.StorageURL, r.URL.String()[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// обработка POST запроса в формате JSON
func (h *AppHandler) CreateHashURLJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// преобразуем bytes (JSON) в map
	var urlJSON models.OriginURL
	if err := json.Unmarshal(data, &urlJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !utils.CheckURL(urlJSON.URL) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	urlHash, err := core.CreateShotURL(&core.InputValueForWriteFile{
		Stor:        h.a.StorageURL,
		BaseAddress: h.a.BaseAddress,
		URL:         urlJSON.URL,
	})

	if err != nil {
		// ошибка при генерации сокращенного URL, возращаем 500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(models.ResponseHashURL{URL: urlHash})
	if err != nil {
		// ошибка при сериализации JSON объекта
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		// ошибка при записи ответа в Body, возращаем 500
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// обработка GET запроса для ping DataBase
func (h *AppHandler) PingDB(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := core.PingDB(h.a.DSNdatabase)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// обработка POST запроса для создания сокращенных url для списка url
func (h *AppHandler) CreateHashURLBatchJSON(_ http.ResponseWriter, _ *http.Request) {
}
