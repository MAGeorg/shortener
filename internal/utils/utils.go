package utils

import (
	"hash/fnv"
	"net/url"
)

// проверка корректности адреса.
func CheckURL(s string) bool {
	u, err := url.Parse(s)
	return err != nil || u.Host != "" || u.Scheme != ""
}

// получение хэша от строки.
func GetHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
