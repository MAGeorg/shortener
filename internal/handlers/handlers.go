package handlers

import (
	"io"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appcontext"
	"github.com/MAGeorg/shortener.git/internal/utils"
)

type AppHandler struct {
	ctx *appcontext.AppContext
	f   func(ctx *appcontext.AppContext, w http.ResponseWriter, r *http.Request)
}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f(h.ctx, w, r)
}

// обработка POST запроса
func CreateHashURL(ctx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	defer r.Body.Close()
	urlStr, err := io.ReadAll(r.Body)

	if err != nil || !utils.CheckURL(string(urlStr)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	urlHash, err := ctx.StorageURL.AddURL(ctx.Cfg.BaseAddress, string(urlStr))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(urlHash))
}

// обработка GET запросв
func GetOriginURL(ctx *appcontext.AppContext, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if len(r.URL.String()) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := ctx.StorageURL.GetOriginURL(r.URL.String()[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
