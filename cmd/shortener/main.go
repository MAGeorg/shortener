package main

import (
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", handlers.CreateHashURL)
	r.Get("/{id}", handlers.GetOriginURL)

	if err := http.ListenAndServe(`:8080`, r); err != nil {
		panic(err)
	}
}
