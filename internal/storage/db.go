// пакет для работы с БД в качестве хранилища
package storage

import (
	"context"
	"database/sql"
	"strconv"
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
func (s *StorageURLinDB) CreateShotURL(url string, h uint32) (string, error) {
	_, err := s.conn.ExecContext(context.Background(),
		"INSERT INTO shot_url (hash_value, origin_url) VALUES ($1,$2);", h, url)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(h), 10), nil
}

// получение из БД изначального запроса по hash
func (s *StorageURLinDB) GetOriginURL(str string) (string, error) {
	res, err := s.conn.QueryContext(context.Background(),
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
