// пакет для работы с БД в качестве хранилища
package storage

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/MAGeorg/shortener.git/internal/models"
)

// структура хранилища БД
//
//nolint:revive // FP
type StorageURLinDB struct {
	conn *sql.DB
}

// создание нового экземпляра хранилища
func NewStorageURLinDB(c *sql.DB) *StorageURLinDB {
	return &StorageURLinDB{
		conn: c,
	}
}

// создание записи в БД с новым сокращенным URL
func (s *StorageURLinDB) CreateShotURL(ctx context.Context, url string, h uint32) (string, error) {
	_, err := s.conn.ExecContext(ctx,
		"INSERT INTO shot_url (hash_value, origin_url) VALUES ($1,$2);", h, url)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(h), 10), nil
}

// получение из БД изначального запроса по hash
func (s *StorageURLinDB) GetOriginURL(ctx context.Context, str string) (string, error) {
	res, err := s.conn.QueryContext(ctx,
		"SELECT origin_url FROM shot_url WHERE hash_value = $1;", str)
	if err != nil || res.Err() != nil {
		return "", err
	}

	defer res.Close()

	var url string
	for res.Next() {
		err = res.Scan(&url)
		if err != nil {
			return "", err
		}
	}

	return url, nil
}

// добавление в БД значений пачкой
func (s *StorageURLinDB) CreateShotURLBatch(context.Context, []models.DataBatch) error {
	return nil
}
