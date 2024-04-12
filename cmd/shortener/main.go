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
	// lastID одновременно служит последним ID записи в файле и флагом для
	// отмены записи в файл при значении -1
	lastID := -1
	if cfg.StorageFileName != "" {
		lastID, err = storage.RestoreData(cfg.StorageFileName, storURL)
		if err != nil {
			logger.Sugar.Errorln("error restore data", err.Error())
			return
		}
	}

	// объявление Producer для записи данных в файл
	var producer *storage.Producer = nil
	if lastID != -1 {
		producer, err = storage.NewProducer(cfg.StorageFileName)
		if err != nil {
			logger.Sugar.Errorln("error get producer", err.Error())
			return
		}
	}

	// инициализация контекста
	appData := appdata.NewAppData(cfg.BaseAddress, storURL, lastID, lg, producer)

	// запуск сервера
	err = handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
