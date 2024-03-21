package handlers

import (
	"log"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/go-chi/chi/v5"
)

func RunServer(a *appdata.AppData) error {
	h := AppHandler{a}

	r := chi.NewRouter()
	r.Post("/", h.CreateHashURL)
	r.Get("/{id}", h.GetOriginURL)

	log.Printf("Server fun on %s address ...", a.Cfg.Address)
	if err := http.ListenAndServe(a.Cfg.Address, r); err != nil {
		return err
	}
	return nil
}
