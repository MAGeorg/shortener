package main

import (
	"github.com/MAGeorg/shortener.git/internal/appcontext"
	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/handlers"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

func main() {
	// парсинг конфига
	cfg := config.NewConfig()
	config.Parse(cfg)

	// инициализация хранилища
	storURL := storage.NewStorageURL()

	// инициализация контекста
	appContext := appcontext.NewAppContext(*cfg, storURL)

	// запуск сервера
	handlers.RunServer(appContext)
}
