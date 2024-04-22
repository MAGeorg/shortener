// пакет appdata реализовывает необходимые данные, которые затем встраиваются в
// в handler.
package appdata

import (
	"github.com/MAGeorg/shortener.git/internal/storage"
	"go.uber.org/zap"
)

// структура AppData содержит base адрес для сокрещенного url
// хранилище, логгер.
type AppData struct {
	BaseAddress string
	StorageURL  storage.Storage
	DSNdatabase string
	Logger      *zap.SugaredLogger
}

// возвращает новый экземпляр AppData.
func NewAppData(baseAddress string, strg storage.Storage, d string, lg *zap.SugaredLogger) *AppData {
	return &AppData{
		BaseAddress: baseAddress,
		StorageURL:  strg,
		DSNdatabase: d,
		Logger:      lg,
	}
}
