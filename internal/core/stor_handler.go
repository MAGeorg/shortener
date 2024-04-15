package core

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

type InputValueForWriteFile struct {
	Stor        *storage.StorageURL
	Producer    *storage.Producer
	BaseAddress string
	URL         string
	LastID      *int
}

// функция реализует бизнес логику обработки начального URL
func CreateShotURL(i *InputValueForWriteFile) (string, error) {
	urlHash, hash, err := i.Stor.AddURL(i.BaseAddress, i.URL)
	if err != nil {
		return "", fmt.Errorf("error add url to storage")
	}
	if *i.LastID != -1 {
		err := i.Producer.WriteEvent(&models.Event{ID: *i.LastID, HashURL: hash, URL: i.URL})
		if err != nil {
			return "", fmt.Errorf("error write value in file")
		}
		*i.LastID += 1
	}
	return urlHash, nil
}

// функция реализует бизнес логику получения начального URL
func GetOriginURL(stor *storage.StorageURL, hash string) (string, error) {
	url, err := stor.GetOriginURL(hash)
	if err != nil {
		return "", fmt.Errorf("error get value from storage")
	}
	return url, nil
}

// функция ping DB
func PingDB(dsn string) error {
	conn, err := ConnectDB(dsn)
	if err != nil {
		return nil
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

// функция реализует бизнес логику подключения к DB
func ConnectDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
