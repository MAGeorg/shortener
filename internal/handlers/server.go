package handlers

import (
	"log"
	"net/http"

	"github.com/MAGeorg/shortener.git/internal/appcontext"
	"github.com/go-chi/chi/v5"
)

func RunServer(ctx *appcontext.AppContext) {
	r := chi.NewRouter()

	r.Method("POST", "/", AppHandler{ctx, CreateHashURL})
	r.Method("GET", "/{id}", AppHandler{ctx, GetOriginURL})

	log.Printf("Server fun on %s address ...", ctx.Cfg.Address)
	if err := http.ListenAndServe(ctx.Cfg.Address, r); err != nil {
		panic(err)
	}
}
