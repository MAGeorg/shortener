package main

import (
	"github.com/MAGeorg/shortener.git/internal/appdata"
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
	appData := appdata.NewAppData(*cfg, storURL)

	// запуск сервера
	err := handlers.RunServer(appData)
	if err != nil {
		panic(err)
	}
}
