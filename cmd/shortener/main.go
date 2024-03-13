package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var savedURL = map[uint32]string{}

// получение хэша от строки
func getHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func checkURL(s string) bool {
	u, err := url.Parse(s)
	return err != nil || u.Host != "" || u.Scheme != ""
}

// генерация сокращенного URL
func generateURL(s string) uint32 {
	h := getHash(s)
	savedURL[h] = s
	return h
}

// получение изначального URL по хэшу
func getOriginURL(s string) (string, error) {
	u := strings.Split(s, "/")
	if len(u) < 1 {
		return "", fmt.Errorf("empty shortner url")
	}

	urlHash, err := strconv.ParseUint(u[len(u)-1], 10, 32)
	if err != nil {
		return "", err
	}

	if k, ok := savedURL[uint32(urlHash)]; ok {
		return k, nil
	}

	return "", fmt.Errorf("not found hash")
}

// CommonHandler обрабатывает запрос в зависимости от метода
func CommonHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		urlStr, err := io.ReadAll(r.Body)

		if err != nil || !checkURL(string(urlStr)) {
			w.WriteHeader(http.StatusNotFound)
			break
		}

		urlHash := fmt.Sprintf("http://%s/%s", r.Host, strconv.FormatUint(uint64(generateURL(string(urlStr))), 10))

		// формирование положительного ответа
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(urlHash))

	case http.MethodGet:
		u := r.URL.String()[1:]

		urlOrig, err := getOriginURL(u)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			break
		}

		// формирование положительного ответа
		w.Header().Set("Location", urlOrig)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", CommonHandle)

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}
