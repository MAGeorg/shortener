package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
	"go.uber.org/zap"
)

type AppData struct {
	BaseAddress string
	StorageURL  *storage.StorageURL
	LastID      int
	Logger      *zap.SugaredLogger
	*storage.Producer
}

func NewAppData(baseAddress string, strg *storage.StorageURL, id int, lg *zap.SugaredLogger, s *storage.Producer) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
		LastID:      id,
		Logger:      lg,
		Producer:    s,
	}
}
