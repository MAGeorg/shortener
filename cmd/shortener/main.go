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

	// инициализация контекста
	appData := appdata.NewAppData(cfg.BaseAddress, nil, cfg.PostgreSQLDSN, lg)

	// инициализация хранилища в зависимости от типа хранилища
	if cfg.StorageFileName != "" {
		producer, err := storage.NewProducer(cfg.StorageFileName)
		if err != nil {
			logger.Sugar.Errorln("error get producer", err.Error())
			return
		}

		storURL := storage.NewStorageURLinFile(producer)
		err = storURL.RestoreData(cfg.StorageFileName)

		if err != nil {
			logger.Sugar.Errorln("error restore data", err.Error())
			return
		}
		appData.StorageURL = storURL
	} else {
		storURL := storage.NewStorageURLinMemory()
		appData.StorageURL = storURL

	}

	// запуск сервера
	err = handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
