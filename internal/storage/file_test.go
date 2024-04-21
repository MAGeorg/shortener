package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/MAGeorg/shortener.git/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	// данные для тестов
	tests := []models.Event{
		{
			ID:      0,
			HashURL: 759827921,
			URL:     "http://yandex.ru",
		},
		{
			ID:      1,
			HashURL: 759827922,
			URL:     "http://ya.ru",
		},
		{
			ID:      2,
			HashURL: 759827923,
			URL:     "https://ex.ru",
		},
	}

	asserts := assert.New(t)

	path := "/tmp/short-url-db.json"
	defer func() {
		err := os.Remove(path)
		asserts.Empty(err)
	}()

	producer, err := NewProducer(path)
	asserts.Empty(err)
	defer producer.Close()

	for i, test := range tests {
		t.Run(fmt.Sprintf("test-producer-%d", i), func(t *testing.T) {
			err := producer.WriteEvent(&test)
			asserts.Empty(err)
		})
	}
}

func TestNewConsumer(t *testing.T) {
	// данные для тестов
	tests := []struct {
		name  string
		write models.Event
		want  models.Event
	}{
		{
			name: "test-1",
			write: models.Event{
				ID:      0,
				HashURL: 759827921,
				URL:     "http://yandex.ru",
			},
			want: models.Event{
				ID:      0,
				HashURL: 759827921,
				URL:     "http://yandex.ru",
			},
		},
		{
			name: "test-2",
			write: models.Event{
				ID:      1,
				HashURL: 759827922,
				URL:     "http://ya.ru",
			},
			want: models.Event{
				ID:      1,
				HashURL: 759827922,
				URL:     "http://ya.ru",
			},
		},
	}

	asserts := assert.New(t)

	path := "/tmp/short-url-db.json"
	defer func() {
		err := os.Remove(path)
		asserts.Empty(err)
	}()

	// наполняем файл данными
	producer, err := NewProducer(path)
	asserts.Empty(err)
	for _, test := range tests {
		err := producer.WriteEvent(&test.write)
		asserts.Empty(err)
	}
	err = producer.Close()
	asserts.Empty(err)

	// считываем и сверяем данные
	consumer, err := NewConsumer(path)
	asserts.Empty(err)
	defer func() {
		err := consumer.Close()
		asserts.Empty(err)
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := consumer.ReadEvent()

			asserts.Empty(err)
			asserts.Equal(test.want.ID, res.ID)
			asserts.Equal(test.want.HashURL, res.HashURL)
			asserts.Equal(test.want.URL, res.URL)

		})
	}

}

func TestRestoreData(t *testing.T) {
	dataWrite := []models.Event{
		{
			ID:      0,
			HashURL: 759827921,
			URL:     "http://yandex.ru",
		},
		{
			ID:      1,
			HashURL: 759827922,
			URL:     "http://ya.ru",
		},
		{
			ID:      2,
			HashURL: 759827923,
			URL:     "https://ex.ru",
		},
	}

	want := []struct {
		lastID  int
		hash    string
		wantURL string
	}{
		{
			lastID:  2,
			hash:    "759827921",
			wantURL: "http://yandex.ru",
		},
		{
			lastID:  2,
			hash:    "759827922",
			wantURL: "http://ya.ru",
		},
		{
			lastID:  2,
			hash:    "759827923",
			wantURL: "https://ex.ru",
		},
	}

	asserts := assert.New(t)

	path := "/tmp/short-url-db.json"
	defer func() {
		err := os.Remove(path)
		asserts.Empty(err)
	}()

	// заполнение файла
	producer, err := NewProducer(path)
	asserts.Empty(err)

	for _, data := range dataWrite {
		err := producer.WriteEvent(&data)
		asserts.Empty(err)
	}
	err = producer.Close()
	asserts.Empty(err)

	// инициализация хранилища
	producer, err = NewProducer(path)
	asserts.Empty(err)
	storURL := NewStorageURLinFile(producer)

	t.Run("test restore data 1", func(t *testing.T) {
		err := storURL.RestoreData(path)
		asserts.Empty(err)

		for _, w := range want {
			v, err := storURL.GetOriginURL(w.hash)
			asserts.Empty(err)
			asserts.Equal(w.wantURL, v)
		}

	})

}
