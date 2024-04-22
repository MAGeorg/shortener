package main

import (
	"context"

	"github.com/MAGeorg/shortener.git/internal/appdata"
	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/core"
	"github.com/MAGeorg/shortener.git/internal/handlers"
	"github.com/MAGeorg/shortener.git/internal/logger"
	"github.com/MAGeorg/shortener.git/internal/storage"
	"github.com/MAGeorg/shortener.git/internal/storage/migration"
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
	switch {
	case cfg.PostgreSQLDSN != "":
		// создаем соединение
		conn, err := core.ConnectDB(cfg.PostgreSQLDSN)
		if err != nil {
			logger.Sugar.Errorln("error connect to db", err.Error())
			return
		}

		// проверяем доступна ли база
		if err := conn.Ping(); err != nil {
			logger.Sugar.Errorln("error open connect: ", err.Error())
			return
		}

		// проверяем есть ли схема
		if res := migration.CheckExistScheme(context.Background(), conn); !res {
			// выполняем миграцию если схемы нет
			source := "../../internal/storage/migration/postgres/001.init_schema.sql"
			migrate := migration.Migration{Source: source}
			err = migrate.Up(conn)
			if err != nil {
				logger.Sugar.Errorln("error execute migrate")
			}
		}

		storURL := storage.NewStorageURLinDB(conn)
		appData.StorageURL = storURL

	case cfg.StorageFileName != "":
		producer, err := storage.NewProducer(cfg.StorageFileName)
		if err != nil {
			logger.Sugar.Errorln("error get producer", err.Error())
			return
		}

		storURL := storage.NewStorageURLinFile(producer)
		// восстанавливаем данные из файла
		err = storURL.RestoreData(cfg.StorageFileName)

		if err != nil {
			logger.Sugar.Errorln("error restore data", err.Error())
			return
		}
		appData.StorageURL = storURL
	default:
		storURL := storage.NewStorageURLinMemory()
		appData.StorageURL = storURL
	}

	// запуск сервера
	err = handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
