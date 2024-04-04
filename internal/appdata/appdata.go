package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
)

type AppData struct {
	BaseAddress string
	StorageURL  *storage.StorageURL
	LastID      int
	*storage.Producer
}

func NewAppData(baseAddress string, strg *storage.StorageURL, id int, s *storage.Producer) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
		LastID:      id,
		Producer:    s,
	}
}
