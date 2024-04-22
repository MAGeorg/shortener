// функионал, описывающий бизнес логику.
package core

import (
	"database/sql"
	"fmt"

	// подключение драйвера PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MAGeorg/shortener.git/internal/storage"
	"github.com/MAGeorg/shortener.git/internal/utils"
)

// структура, содержащая необходимы параметры для создания нового соркащенного URL.
type InputValueForWriteFile struct {
	Stor        storage.Storage
	BaseAddress string
	URL         string
}

// функция реализует бизнес логику обработки начального URL.
func CreateShotURL(i *InputValueForWriteFile) (string, error) {
	if !utils.CheckURL(i.URL) {
		return "", fmt.Errorf("not valid url")
	}

	h := utils.GetHash(i.URL)
	urlHash, err := i.Stor.CreateShotURL(i.URL, h)

	if err != nil {
		return "", fmt.Errorf("error add url to storage: %w", err)
	}
	return fmt.Sprintf("%s/%s", i.BaseAddress, urlHash), nil
}

// функция реализует бизнес логику получения начального URL.
func GetOriginURL(stor storage.Storage, hash string) (string, error) {
	url, err := stor.GetOriginURL(hash)
	if err != nil {
		return "", fmt.Errorf("error get value from storage: %w", err)
	}
	return url, nil
}

// функция ping DB.
func PingDB(dsn string) error {
	conn, err := ConnectDB(dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

// функция реализует бизнес логику подключения к DB.
func ConnectDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
