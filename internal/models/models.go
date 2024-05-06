// пакет, содержащий описание используемых моделей для запросов/ответов.
package models

// структура для обработки POST запроса, содержащего
// URL в теле в формате JSON.
type OriginURL struct {
	URL string `json:"url"`
}

// структура для обработки ответа POST запроса, содержащего
// сокращенный URL в теле в формате JSON.
type ResponseHashURL struct {
	URL string `json:"result"`
}

// структура для записи/считывания из файла.
type Event struct {
	URL     string `json:"original_url"`
	ID      int    `json:"id"`
	HashURL uint32 `json:"short_url"`
}

// структура для обработки тела запроса батчами.
type DataBatch struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginURL     string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	Hash          uint32 `json:"-"`
}
