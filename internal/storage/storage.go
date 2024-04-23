package storage

import (
	"context"

	"github.com/MAGeorg/shortener.git/internal/models"
)

// интерфейс, описывающий контракты для хранилищ (БД, файл, память)
type Storage interface {
	// input (context, originURL, hash) (shotHashURL, error)
	CreateShotURL(context.Context, string, uint32) (string, error)
	// input (hash) (originURL, error)
	GetOriginURL(context.Context, string) (string, error)
	// input (context, []models.DataBatch)
	CreateShotURLBatch(context.Context, []models.DataBatch) error
}
