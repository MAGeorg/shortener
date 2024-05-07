package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid url",
			url:  "https://practicum.yandex.ru/",
			want: true,
		},
		{
			name: "valid detail url",
			url:  "http://127.0.0.1:8000/home/",
			want: true,
		},
		{
			name: "incorrect url",
			url:  "http",
			want: false,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			ans := CheckURL(test.url)

			asserts.Equal(test.want, ans)
		})
	}
}
