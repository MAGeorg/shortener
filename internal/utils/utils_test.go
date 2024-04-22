package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want uint32
	}{
		{
			name: "valid string",
			s:    "http://hello.ru",
			want: uint32(1713553680),
		},
		{
			name: "empty string",
			s:    "",
			want: uint32(2166136261),
		},
		{
			name: "numbers string",
			s:    "0123456789",
			want: uint32(4185952242),
		},
		{
			name: "special character string",
			s:    "=-+)(')'\n",
			want: uint32(1995809324),
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(_ *testing.T) {
			ans := GetHash(test.s)

			asserts.Equal(test.want, ans)
		})
	}
}

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
