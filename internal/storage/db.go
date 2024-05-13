// пакет для работы с БД в качестве хранилища.
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	customerr "github.com/MAGeorg/shortener.git/internal/errors"
	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

// структура хранилища БД.
//
//nolint:revive // FP
type StorageURLinDB struct {
	conn *sql.DB
}

// создание нового экземпляра хранилища.
func NewStorageURLinDB(c *sql.DB) *StorageURLinDB {
	return &StorageURLinDB{
		conn: c,
	}
}

// создание записи в БД с новым сокращенным URL.
func (s *StorageURLinDB) CreateShotURL(ctx context.Context, url string, h uint32, userID int) (string, error) {
	_, err := s.conn.ExecContext(ctx,
		"INSERT INTO shot_url (hash_value, origin_url, user_id) VALUES ($1,$2,$3);", h, url, userID)

	// проверка на дубликат
	if driverErr, ok := err.(*pgconn.PgError); ok && driverErr.Code == "23505" {
		return strconv.FormatUint(uint64(h), 10), customerr.ErrAccessDenied
	}
	return strconv.FormatUint(uint64(h), 10), err
}

// получение из БД изначального запроса по hash.
func (s *StorageURLinDB) GetOriginURL(ctx context.Context, str string, _ int) (string, error) {
	res, err := s.conn.QueryContext(ctx,
		"SELECT origin_url, is_deleted FROM shot_url WHERE hash_value = $1;", str)
	if err != nil || res.Err() != nil {
		return "", err
	}

	defer res.Close()

	var (
		url string
		del sql.NullBool
	)

	for res.Next() {
		err = res.Scan(&url, &del)
		if err != nil {
			return "", err
		}
	}

	if del.Valid && del.Bool {
		return url, customerr.ErrDeleteShotURL
	}
	return url, nil
}

// добавление в БД значений пачкой.
func (s *StorageURLinDB) CreateShotURLBatch(ctx context.Context, d []models.DataBatch, userID int) error {
	// начинаем транзакцию.
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}

	// выполняем запись.
	for _, i := range d {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO shot_url (hash_value, origin_url, user_id) VALUES ($1,$2,$3);",
			i.Hash, i.OriginURL, userID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// коммитим изменения.
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// получение всех пар short_url - original_url.
func (s *StorageURLinDB) GetAllURL(ctx context.Context, baseAddr string, userID int) ([]models.DataBatch, error) {
	res, err := s.conn.QueryContext(ctx,
		"SELECT hash_value, origin_url FROM shot_url WHERE user_id = $1", userID)
	if err != nil || res.Err() != nil {
		return nil, err
	}

	defer res.Close()

	r := []models.DataBatch{}
	var (
		hash    uint32
		origURL string
	)

	for res.Next() {
		err := res.Scan(&hash, &origURL)

		if err != nil {
			return nil, fmt.Errorf("error scan value from db: %w", err)
		}

		r = append(r, models.DataBatch{
			ShortURL:  fmt.Sprintf("%s/%s", baseAddr, strconv.FormatUint(uint64(hash), 10)),
			OriginURL: origURL,
		})
	}

	return r, nil
}

// удаление по hash.
func (s *StorageURLinDB) DeleteValueByHash(ctx context.Context, hash uint32, userID int) error {
	_, err := s.conn.ExecContext(ctx,
		"UPDATE shot_url SET is_deleted = TRUE WHERE hash_value = $1 AND user_id = $2;",
		hash, userID)
	return err
}
