package main

import (
	"log"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", handlers.CreateHashURL)
	r.Get("/{id}", handlers.GetOriginURL)

	log.Printf("Server fun on %s address ...", config.Conf.Address)
	if err := http.ListenAndServe(config.Conf.Address, r); err != nil {
		panic(err)
	}
}
