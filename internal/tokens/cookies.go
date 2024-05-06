// содержит функционал для получения информации из куков и установки куков.
package tokens

import (
	"fmt"
	"net/http"
)

// функция GetValueFromCookie выполняет поиск в cookie по ключу и возвращает значение.
func GetValueFromCookie(r *http.Request, key string) string {
	c := r.Cookies()
	for _, i := range c {
		if i.Name == key {
			return i.Value
		}
	}
	return ""
}

// функция SetValueToCookie выставляет значение в cookie.
func SetValueToCookie(w http.ResponseWriter, values ...string) error {
	if len(values)%2 != 0 {
		return fmt.Errorf("incomplete key-value pair")
	}
	for i := 0; i < len(values); i += 2 {
		http.SetCookie(w, &http.Cookie{Name: values[i], Value: values[i+1]})
	}
	return nil
}
