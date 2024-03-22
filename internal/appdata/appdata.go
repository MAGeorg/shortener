package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

type AppData struct {
	Cfg        config.Config
	StorageURL *storage.StorageURL
}

func NewAppData(cfg config.Config, strg *storage.StorageURL) *AppData {
	return &AppData{
		Cfg:        cfg,
		StorageURL: strg,
	}
}
