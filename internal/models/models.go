// пакет, содержащий описание используемых моделей для запросов/ответов
package models

// структура для обработки POST запроса, содержащего
// URL в теле в формате JSON
type OriginURL struct {
	URL string `json:"url"`
}

// структура для обработки ответа POST запроса, содержащего
// сокращенный URL в теле в формате JSON
type ResponseHashURL struct {
	URL string `json:"result"`
}

// структура для записи/считывания из файла
type Event struct {
	ID      int    `json:"id"`
	HashURL uint32 `json:"short_url"`
	URL     string `json:"original_url"`
}
