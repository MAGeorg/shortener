// пакет для работы с файлом в качестве хранилища
package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MAGeorg/shortener.git/internal/models"
)

// структура хранилища файла
//
//nolint:revive // FP
type StorageURLinFile struct {
	Producer *Producer
	savedURL map[uint32]string
	lastID   int
}

// получение нового экземпляра хранилища URL по хэшу
func NewStorageURLinFile(s *Producer) *StorageURLinFile {
	return &StorageURLinFile{
		Producer: s,
		savedURL: make(map[uint32]string),
		lastID:   0,
	}
}

// создание записи в файле с новым сокращенным URL
func (s *StorageURLinFile) CreateShotURL(ctx context.Context, url string, h uint32) (string, error) {
	// проверяем, есть ли уже запись в файле и локальном кэше
	if _, ok := s.savedURL[h]; ok {
		return strconv.FormatUint(uint64(h), 10), nil
	}

	err := s.Producer.WriteEvent(&models.Event{ID: s.lastID, HashURL: h, URL: url})
	if err != nil {
		return "", fmt.Errorf("error write value in file")
	}
	s.savedURL[h] = url
	s.lastID++
	return strconv.FormatUint(uint64(h), 10), nil
}

// получение из БД изначального запроса по hash
func (s *StorageURLinFile) GetOriginURL(ctx context.Context, str string) (string, error) {
	// преобразование строки с HashURL в uint32
	urlHash, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return "", fmt.Errorf("incorrect hash")
	}

	// поиск оригинального адреса по HashURL
	urlOrig, ok := s.savedURL[uint32(urlHash)]
	if !ok {
		return "", fmt.Errorf("not found url by hash")
	}
	return urlOrig, nil
}

// функция для добавления в файл данных пачкой
func (s *StorageURLinFile) CreateShotURLBatch(context.Context, []models.DataBatch) error {
	return nil
}

// функция восстановления данных и записи в хранилище в памяти
func (s *StorageURLinFile) RestoreData(path string) error {
	var lastID int
	consumer, err := NewConsumer(path)
	if err != nil {
		return err
	}
	defer consumer.Close()

	for {
		e, err := consumer.ReadEvent()
		if err != nil || e == nil {
			break
		}
		lastID = e.ID
		s.savedURL[e.HashURL] = e.URL
	}
	s.lastID = lastID
	//nolint:nilerr // на 74 строке ошибка считывания EOF
	return nil
}

// структура Consumer, содержит указатель на файл, с которым работаем
// и scanner
type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

// создание экземпляра Consumer
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// запись consumer события (новой записи сокращенного URL)
func (c *Consumer) ReadEvent() (*models.Event, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()

	event := models.Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// закрывает файл, с которым работает consumer
func (c *Consumer) Close() error {
	return c.file.Close()
}

// структура Producer, содержит указать на файл, с которым работает
// и writer
type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

// получение нового экземпляра Producer
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

// запись
func (p *Producer) WriteEvent(event *models.Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

// закрытие файла, с которым работает producer
func (p *Producer) Close() error {
	return p.file.Close()
}
