package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) Save(filename string, content io.Reader) (string, error) {
	path := filepath.Join(s.basePath, filename)
	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, content)
	if err != nil {
		return "", fmt.Errorf("failed to write file content: %w", err)
	}

	return path, nil
}

func (s *LocalStorage) Get(filename string) (string, error) {
	path := filepath.Join(s.basePath, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %w", err)
	}
	return path, nil
}

func (s *LocalStorage) Delete(filename string) error {
	path := filepath.Join(s.basePath, filename)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
