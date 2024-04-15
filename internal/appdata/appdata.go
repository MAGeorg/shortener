package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
	"go.uber.org/zap"
)

type AppData struct {
	BaseAddress string
	StorageURL  *storage.StorageURL
	DSNdatabase string
	LastID      int
	Logger      *zap.SugaredLogger
	*storage.Producer
}

func NewAppData(baseAddress string, strg *storage.StorageURL, d string, id int, lg *zap.SugaredLogger, s *storage.Producer) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
		DSNdatabase: d,
		LastID:      id,
		Logger:      lg,
		Producer:    s,
	}
}
