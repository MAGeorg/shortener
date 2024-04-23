// пакет для работы с БД в качестве хранилища
package storage

import (
	"context"
	"database/sql"
	"strconv"

	customerr "github.com/MAGeorg/shortener.git/internal/errors"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
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

	// проверка на дубликат
	if driverErr, ok := err.(*pgconn.PgError); ok && driverErr.Code == "23505" {
		return strconv.FormatUint(uint64(h), 10), customerr.ErrAccessDenied
	}
	return strconv.FormatUint(uint64(h), 10), err
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
func (s *StorageURLinDB) CreateShotURLBatch(ctx context.Context, d []models.DataBatch) error {
	// начинаем транзакцию
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	// выполняем запись
	for _, i := range d {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO shot_url (hash_value, origin_url) VALUES ($1,$2);", i.Hash, i.OriginURL)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// коммитим изменения
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
