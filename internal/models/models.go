package models

// структура для обработки запроса POST запроса, содержащего
// URL в теле в формате JSON
type OriginURL struct {
	URL string `json:"url"`
}

type AnswerHashURL struct {
	URL string `json:"result"`
}
