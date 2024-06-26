package handlers

import (
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/logger"
	"github.com/MAGeorg/shortener.git/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// функция инициализации endpoint и запуска сервера.
func RunServer(address string, a *appdata.AppData) error {
	h := AppHandler{a}
	lgMiddleware := logger.NewLogMiddleware(a.Logger)

	r := chi.NewRouter()

	r.Method("POST", "/", lgMiddleware.LogMiddleware(middleware.GzipMiddleware(http.HandlerFunc(h.CreateHashURL))))
	r.Method("POST", "/api/shorten",
		lgMiddleware.LogMiddleware(middleware.GzipMiddleware(http.HandlerFunc(h.CreateHashURLJSON))))
	r.Method("POST", "/api/shorten/batch",
		lgMiddleware.LogMiddleware(middleware.GzipMiddleware(http.HandlerFunc(h.CreateHashURLBatchJSON))))
	r.Method("GET", "/{id}", lgMiddleware.LogMiddleware(middleware.GzipMiddleware(http.HandlerFunc(h.GetOriginURL))))
	r.Method("GET", "/ping", lgMiddleware.LogMiddleware(middleware.GzipMiddleware(http.HandlerFunc(h.PingDB))))

	a.Logger.Infof("Server run on %s address ...", address)
	//nolint:gosec // no matter in this
	if err := http.ListenAndServe(address, r); err != nil {
		return err
	}
	return nil
}
