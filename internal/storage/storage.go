package storage

import (
	"fmt"
	"strconv"

	"github.com/MAGeorg/shortener.git/internal/utils"
)

// структура хранилища URL по хэшу
type StorageURL struct {
	savedURL map[uint32]string
}

// получение нового экземпляра хранилища URL по хэшу
func NewStorageURL() *StorageURL {
	return &StorageURL{
		savedURL: make(map[uint32]string),
	}
}

func (s *StorageURL) AddURL(base, url string) (string, error) {
	// проверка валидности URL
	if utils.CheckURL(url) {
		h := utils.GetHash(url)
		s.savedURL[h] = url
		return fmt.Sprintf("%s/%s", base, strconv.FormatUint(uint64(h), 10)), nil
	} else {
		return "", fmt.Errorf("not valid url")
	}
}

func (s *StorageURL) GetOriginURL(str string) (string, error) {
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
