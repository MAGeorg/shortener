package storage

import (
	"database/sql"
)

type StorageURLinDB struct {
	conn *sql.DB
}

func NewStorageURLinDB(c *sql.DB) *StorageURLinDB {
	return &StorageURLinDB{
		conn: c,
	}
}

func (s *StorageURLinDB) CreateShotURL(url string, h uint32) (string, error) {
	return "", nil
}

func (s *StorageURLinDB) GetOriginURL(str string) (string, error) {
	return "", nil
}
