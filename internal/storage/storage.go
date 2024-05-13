package storage

import (
	"context"

	"github.com/MAGeorg/shortener.git/internal/models"
)

// интерфейс, описывающий контракты для хранилищ (БД, файл, память).
type Storage interface {
	// input (context, originURL, hash, userID) (shotHashURL, error).
	CreateShotURL(context.Context, string, uint32, int) (string, error)
	// input (context, hash, userID) (originURL, error).
	GetOriginURL(context.Context, string, int) (string, error)
	// input (context, []models.DataBatch, userID).
	CreateShotURLBatch(context.Context, []models.DataBatch, int) error
	// input (context, baseAddress, userID).
	GetAllURL(context.Context, string, int) ([]models.DataBatch, error)
}
