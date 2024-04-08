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
	lg, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := lg.Sync(); err != nil {
			panic(err)
		}
	}()

	// инициализация хранилища
	storURL := storage.NewStorageURL()

	// проверка конфига на наличие переменной для хранения имени файла
	var lastID int
	if cfg.StorageFileName != "" {
		lastID, err = storage.RestoreData(cfg.StorageFileName, storURL)
		if err != nil {
			logger.Sugar.Errorln("error restore data", err.Error())
			return
		}
	} else {
		// выставляем как флаг, что в файл сохранять данные не нужно
		lastID = -1
	}

	var producer *storage.Producer
	if lastID != -1 {
		producer, err = storage.NewProducer(cfg.StorageFileName)
		if err != nil {
			logger.Sugar.Errorln("error get producer", err.Error())
		}
	} else {
		producer = nil
	}

	// инициализация контекста
	appData := appdata.NewAppData(cfg.BaseAddress, storURL, lastID, lg, producer)

	// запуск сервера
	err = handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
