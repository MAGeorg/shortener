// точка входа для сервера.
package main

import (
	"context"
	"fmt"

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
			//nolint:forbidigo // FP
			fmt.Printf("error sync logger: %s", err.Error())
		}
	}()

	// проверка значений для запуска
	lg.Infof("address = %s base = %s filename = %s postgres-dsn = %s",
		cfg.Address, cfg.BaseAddress, cfg.StorageFileName, cfg.PostgreSQLDSN)

	// инициализация контекста
	appData := appdata.NewAppData(cfg.BaseAddress, nil, cfg.PostgreSQLDSN, lg)

	// инициализация хранилища в зависимости от типа хранилища
	switch {
	case cfg.PostgreSQLDSN != "":
		// создаем соединение
		conn, err := core.ConnectDB(cfg.PostgreSQLDSN)
		if err != nil {
			lg.Errorln("error connect to db", err.Error())
			return
		}

		// проверяем доступна ли база
		if err := conn.Ping(); err != nil {
			lg.Errorln("error open connect: ", err.Error())
			return
		}

		// проверяем есть ли схема
		if res := migration.CheckExistScheme(context.Background(), conn); !res {
			// выполняем миграцию если схемы нет
			source := "../../internal/storage/migration/postgres/001.init_schema.sql"
			migrate := migration.Migration{Source: source}
			err = migrate.Up(conn)
			if err != nil {
				lg.Errorln("error execute migrate", err)
			}
		}

		storURL := storage.NewStorageURLinDB(conn)
		appData.StorageURL = storURL
		lg.Infoln("success add db storage")

	case cfg.StorageFileName != "":
		producer, err := storage.NewProducer(cfg.StorageFileName)
		if err != nil {
			lg.Errorln("error get producer", err.Error())
			return
		}

		storURL := storage.NewStorageURLinFile(producer)
		// восстанавливаем данные из файла
		err = storURL.RestoreData(cfg.StorageFileName)

		if err != nil {
			lg.Errorln("error restore data", err.Error())
			return
		}
		appData.StorageURL = storURL
		lg.Infoln("success add file storage")

	default:
		storURL := storage.NewStorageURLinMemory()
		appData.StorageURL = storURL
		lg.Infoln("success add in-memory storage")
	}

	// запуск сервера
	err = handlers.RunServer(cfg.Address, appData)
	if err != nil {
		panic(err)
	}
}
