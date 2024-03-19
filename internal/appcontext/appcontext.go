package appcontext

import (
	"github.com/MAGeorg/shortener.git/internal/config"
	"github.com/MAGeorg/shortener.git/internal/storage"
)

type AppContext struct {
	Cfg        config.Config
	StorageURL *storage.StorageURL
}

func NewAppContext(cfg config.Config, strg *storage.StorageURL) *AppContext {
	return &AppContext{
		Cfg:        cfg,
		StorageURL: strg,
	}
}
