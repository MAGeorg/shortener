package handlers

import (
	"log"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/logger"
	"github.com/go-chi/chi/v5"
)

func RunServer(address string, a *appdata.AppData) error {
	h := AppHandler{a}

	r := chi.NewRouter()

	r.Method("POST", "/", logger.MiddlewareLog(http.HandlerFunc(h.CreateHashURL)))
	r.Method("POST", "/api/shorten", logger.MiddlewareLog(http.HandlerFunc(h.CreateHashURLJSON)))
	r.Method("GET", "/{id}", logger.MiddlewareLog(http.HandlerFunc(h.GetOriginURL)))

	log.Printf("Server fun on %s address ...", address)
	if err := http.ListenAndServe(address, r); err != nil {
		return err
	}
	return nil
}
