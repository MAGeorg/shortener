package storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MAGeorg/shortener.git/internal/models"
)

// структура хранилища URL по хэшу
//
//nolint:revive // FP
type StorageURLinMemory struct {
	savedURL map[uint32]string
}

// получение нового экземпляра хранилища URL по хэшу
func NewStorageURLinMemory() *StorageURLinMemory {
	return &StorageURLinMemory{
		savedURL: make(map[uint32]string),
	}
}

// создание записи в памяти с новым сокращенным URL
func (s *StorageURLinMemory) CreateShotURL(_ context.Context, url string, h uint32) (string, error) {
	// добавление в хранилище в памяти
	s.savedURL[h] = url
	return strconv.FormatUint(uint64(h), 10), nil
}

// получение из памяти изначального запроса по hash
func (s *StorageURLinMemory) GetOriginURL(_ context.Context, str string) (string, error) {
	// преобразование строки с HashURL в uint32
	urlHash, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return "", fmt.Errorf("incorrect hash")
	}

	// поиск оригинального адреса по HashURL
	urlOrig, ok := s.savedURL[uint32(urlHash)]
	if !ok {
		return "", fmt.Errorf("not found url by hash")
	}
	return urlOrig, nil
}

// добавление в память значений пачкой
func (s *StorageURLinMemory) CreateShotURLBatch(_ context.Context, d []models.DataBatch) error {
	for _, i := range d {
		s.savedURL[i.Hash] = i.OriginURL
	}
	return nil
}
