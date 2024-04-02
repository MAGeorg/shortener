package main

import (
	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/handlers"
	"github.com/MAGeorg/shortener.git/internal/logger"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

func main() {
	// парсинг конфига
	cfg := config.NewConfig()
	config.Parse(cfg)

	// инициализация логгера
	if err := logger.NewLogger(); err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Sugar.Sync(); err != nil {
			panic(err)
		}
	}()

	// инициализация хранилища
	storURL := storage.NewStorageURL()

	// инициализация контекста
	appData := appdata.NewAppData(cfg.BaseAddress, storURL)

	// запуск сервера
	err := handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
