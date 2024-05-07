package handlers

import (
	"net/url"
)

// проверка корректности адреса.
func CheckURL(s string) bool {
	u, err := url.Parse(s)
	return err != nil || u.Host != "" || u.Scheme != ""
}
