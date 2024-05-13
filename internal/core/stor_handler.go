// функионал, описывающий бизнес логику.
package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	// подключение драйвера PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/MAGeorg/shortener.git/internal/hash"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

// структура, содержащая необходимы параметры для создания нового соркащенного URL.
type InputValueForWriteFile struct {
	Stor        storage.Storage
	BaseAddress string
	URL         string
	UserID      int
}

// функция реализует бизнес логику обработки начального URL.
func CreateShotURL(ctx context.Context, i *InputValueForWriteFile) (string, error) {
	h := hash.GetHash(i.URL)
	urlHash, err := i.Stor.CreateShotURL(ctx, i.URL, h, i.UserID)

	if err != nil {
		return fmt.Sprintf("%s/%s", i.BaseAddress, urlHash), fmt.Errorf("error add url to storage: %w", err)
	}
	return fmt.Sprintf("%s/%s", i.BaseAddress, urlHash), nil
}

// функция реализует бизнес логику получения начального URL.
func GetOriginURL(ctx context.Context, stor storage.Storage, hashString string, userID int) (string, error) {
	url, err := stor.GetOriginURL(ctx, hashString, userID)
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

// функция, реализующая бизнес логику обработки batch json.
func CreateShotURLBatch(ctx context.Context, stor storage.Storage,
	base string, d []models.DataBatch, userID int) ([]models.DataBatch, error) {
	res := []models.DataBatch{}

	// заполняем сокращенный url и результат обработки.
	for i := range d {
		d[i].Hash = hash.GetHash(d[i].OriginURL)
		res = append(res, models.DataBatch{
			CorrelationID: d[i].CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%d", base, d[i].Hash),
		})
	}

	err := stor.CreateShotURLBatch(ctx, d, userID)
	if err != nil {
		return res, err
	}

	return res, nil
}

// функция, реализующая бизнес-логику для получения всех значений short_url - original_url.
func GetAllURL(ctx context.Context, stor storage.Storage, base string, userID int) ([]byte, error) {
	res, err := stor.GetAllURL(ctx, base, userID)
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		b, err := json.Marshal(res)
		return b, err
	}
	return nil, fmt.Errorf("empty result")
}
