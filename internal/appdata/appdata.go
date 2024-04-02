package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
)

type AppData struct {
	BaseAddress string
	StorageURL  *storage.StorageURL
}

func NewAppData(baseAddress string, strg *storage.StorageURL) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
	}
}
