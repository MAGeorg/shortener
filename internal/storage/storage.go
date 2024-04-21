package storage

// интерфейс, описывающий контракты для хранилищ (БД, файл, память)
type Storage interface {
	// input (originURL, hash) (shotHashURL, error)
	CreateShotURL(string, uint32) (string, error)
	// input (hash) (originURL, error)
	GetOriginURL(string) (string, error)
}
