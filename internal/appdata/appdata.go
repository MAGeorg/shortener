// пакет appdata реализовывает необходимые данные, которые затем встраиваются в
// в handler.
package appdata

import (
	"go.uber.org/zap"

	"github.com/MAGeorg/shortener.git/internal/storage"
	"github.com/MAGeorg/shortener.git/internal/tokens"
)

// структура AppData содержит base адрес для сокрещенного url
// хранилище, логгер.
//
//nolint:govet // FP
type AppData struct {
	BaseAddress string
	DSNdatabase string
	StorageURL  storage.Storage
	Logger      *zap.SugaredLogger
	Tokens      *tokens.TokensID
}

// возвращает новый экземпляр AppData.
func NewAppData(baseAddr string, s storage.Storage, d string, lg *zap.SugaredLogger, t *tokens.TokensID) *AppData {
	return &AppData{
		BaseAddress: baseAddr,
		StorageURL:  s,
		DSNdatabase: d,
		Logger:      lg,
		Tokens:      t,
	}
}
