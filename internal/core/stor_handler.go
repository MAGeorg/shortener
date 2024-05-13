// функионал, описывающий бизнес логику.
package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	// подключение драйвера PostgreSQL.
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	customerr "github.com/MAGeorg/shortener.git/internal/errors"
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

// структура, содержащая в себе параметры для удаления batch URL.
//
//nolint:govet // FP
type InputValueForDeleteBatch struct {
	Stor   storage.Storage
	Hashs  []string
	UserID int
	Logger *zap.SugaredLogger
}

// функция реализует бизнес-логику обработки начального URL.
func CreateShotURL(ctx context.Context, i *InputValueForWriteFile) (string, error) {
	h := hash.GetHash(i.URL)
	urlHash, err := i.Stor.CreateShotURL(ctx, i.URL, h, i.UserID)

	if err != nil {
		return fmt.Sprintf("%s/%s", i.BaseAddress, urlHash), fmt.Errorf("error add url to storage: %w", err)
	}
	return fmt.Sprintf("%s/%s", i.BaseAddress, urlHash), nil
}

// функция реализует бизнес-логику получения начального URL.
func GetOriginURL(ctx context.Context, stor storage.Storage, hashString string, userID int) (string, error) {
	url, err := stor.GetOriginURL(ctx, hashString, userID)
	return url, err
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

// функция реализует бизнес-логику подключения к DB.
func ConnectDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// функция, реализующая бизнес-логику обработки batch json.
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
	return nil, customerr.ErrEmptyResult
}

// функция для конвертирования string слайса в uint32 слайс.
func convertListStringToListUint32(input []string) []uint32 {
	res := make([]uint32, 0)
	for _, i := range input {
		if u32, err := strconv.ParseUint(i, 10, 32); err == nil {
			res = append(res, uint32(u32))
		}
	}
	return res
}

// фукнкиця, реализующая бизнес-логику для удаления всех значений shot_url по батчу хэшей.
func DeleteBatchURLbyHash(ctx context.Context, input *InputValueForDeleteBatch) error {
	values := convertListStringToListUint32(input.Hashs)

	for _, i := range values {
		go func(hash uint32) {
			err := input.Stor.DeleteValueByHash(ctx, hash, input.UserID)
			if err != nil {
				input.Logger.Infof("error delete %d: %s", hash, err.Error())
			}
		}(i)
	}

	return nil
}
