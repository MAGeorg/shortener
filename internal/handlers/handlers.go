package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MAGeorg/shortener.git/internal/utils"
)

var savedURL = map[uint32]string{}

// генерация сокращенного URL
func generateURL(s string) uint32 {
	h := utils.GetHash(s)
	savedURL[h] = s
	return h
}

// обработка POST запроса
func CreateHashURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	defer r.Body.Close()
	urlStr, err := io.ReadAll(r.Body)

	if err != nil || !utils.CheckURL(string(urlStr)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	urlHash := fmt.Sprintf("http://%s/%s", r.Host, strconv.FormatUint(uint64(generateURL(string(urlStr))), 10))

	// формирование положительного ответа
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(urlHash))
}

// обработка GET запросв
func GetOriginURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	s := r.URL.String()[1:]

	// преобразование строки с HashURL в uint32
	urlHash, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// поиск оригинального адреса по HashURL
	urlOrig, ok := savedURL[uint32(urlHash)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// формирование положительного ответа
	w.Header().Set("Location", urlOrig)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
