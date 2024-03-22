package handlers

import (
	"io"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
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

	urlHash, err := h.a.StorageURL.AddURL(h.a.Cfg.BaseAddress, string(urlStr))

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
