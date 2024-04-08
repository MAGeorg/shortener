package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/MAGeorg/shortener.git/internal/models"
)

type Consumer struct {
	file *os.File
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

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

func (c *Consumer) Close() error {
	return c.file.Close()
}

// функция восстановления данных и записи в хранилище в памяти
// вернет ID последней записи, чтобы продолжить запись, и ошибку
func RestoreData(path string, stor *StorageURL) (int, error) {
	var lastID int
	consumer, err := NewConsumer(path)
	if err != nil {
		return lastID, err
	}
	defer consumer.Close()

	for {
		e, err := consumer.ReadEvent()
		if err != nil || e == nil {
			break
		}
		lastID = e.ID
		stor.Add(e.URL, e.HashURL)
	}

	return lastID, nil
}

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

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

func (p *Producer) WriteEvent(id *int, event *models.Event) error {
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

	*id += 1
	return p.writer.Flush()
}

func (p *Producer) Close() error {
	return p.file.Close()
}
