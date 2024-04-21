package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
	"go.uber.org/zap"
)

type AppData struct {
	BaseAddress string
	StorageURL  storage.Storage
	DSNdatabase string
	Logger      *zap.SugaredLogger
}

func NewAppData(baseAddress string, strg storage.Storage, d string, lg *zap.SugaredLogger) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
		DSNdatabase: d,
		Logger:      lg,
	}
}
