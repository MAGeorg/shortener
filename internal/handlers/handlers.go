package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/MAGeorg/shortener.git/internal/utils"
)

type AppHandler struct {
	a *appdata.AppData
}

// обработка POST запроса
func (h *AppHandler) CreateHashURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	defer r.Body.Close()
	urlStr, err := io.ReadAll(r.Body)

	if err != nil || !utils.CheckURL(string(urlStr)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	urlHash, hash, err := h.a.StorageURL.AddURL(h.a.BaseAddress, string(urlStr))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(urlHash))
	if err != nil {
		// ошибка при записи ответа в Body, возращаем 500
		w.WriteHeader(http.StatusInternalServerError)
	}

	if h.a.LastID != -1 {
		h.a.LastID += 1
		err := h.a.Producer.WriteEvent(&models.Event{ID: h.a.LastID, HashURL: hash, URL: string(urlStr)})
		if err != nil {
			// ошибка при записи в файл, возращаем 500
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// обработка GET запросв
func (h *AppHandler) GetOriginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if len(r.URL.String()) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := h.a.StorageURL.GetOriginURL(r.URL.String()[1:])
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// преобразуем bytes (JSON) в map
	var urlJSON models.OriginURL
	if err := json.Unmarshal(data, &urlJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !utils.CheckURL(urlJSON.URL) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	urlHash, hash, err := h.a.StorageURL.AddURL(h.a.BaseAddress, urlJSON.URL)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(models.AnswerHashURL{URL: urlHash})
	if err != nil {
		// ошибка при сериализации JSON объекта
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		// ошибка при записи ответа в Body, возращаем 500
		w.WriteHeader(http.StatusInternalServerError)
	}
	if h.a.LastID != -1 {
		h.a.LastID += 1
		err := h.a.Producer.WriteEvent(&models.Event{ID: h.a.LastID, HashURL: hash, URL: urlJSON.URL})
		if err != nil {
			// ошибка при записи в файл, возращаем 500
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
