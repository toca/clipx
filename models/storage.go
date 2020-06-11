package models

import (
	"bufio"
	"io"
	"os"
)

type Storage interface {
	Load(path *string) ([]string, error)
	Save(path *string, data []string) error
}

type storage struct {
}

func NewStorage() Storage {
	return &storage{}
}

func (this *storage) Load(path *string) ([]string, error) {
	result := make([]string, 0, 16)
	file, err := os.OpenFile(*path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		buf, err := reader.ReadBytes(byte('\x00'))
		if err != nil && err != io.EOF {
			return nil, err
		}
		if 0 < len(buf) {
			result = append(result, string(buf[:len(buf)-1]))
		}
		if err == io.EOF {
			break
		}
	}
	return result, nil
}

func (this *storage) Save(path *string, data []string) error {
	file, err := os.Create(*path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i, _ := range data {
		_, err := writer.WriteString(data[i])
		if err != nil {
			return err
		}
		writer.WriteRune('\x00')
	}
	writer.Flush()
	return nil
}
