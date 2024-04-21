package storage

import (
	"fmt"
	"strconv"
)

// структура хранилища URL по хэшу
type StorageURLinMemory struct {
	savedURL map[uint32]string
}

// получение нового экземпляра хранилища URL по хэшу
func NewStorageURLinMemory() *StorageURLinMemory {
	return &StorageURLinMemory{
		savedURL: make(map[uint32]string),
	}
}

func (s *StorageURLinMemory) CreateShotURL(url string, h uint32) (string, error) {
	// добавление в хранилище в памяти
	s.savedURL[h] = url
	return strconv.FormatUint(uint64(h), 10), nil
}

func (s *StorageURLinMemory) GetOriginURL(str string) (string, error) {
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
