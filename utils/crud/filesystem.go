package crud

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ErrRecordDoesNotExist represents when file path is not found on file system
var ErrRecordDoesNotExist = errors.New("File does not exist")

// NewFileSystemStore creates a Store backed by a file system directory.
// Each key is represented by a file in that directory.
func NewFileSystemStore(baseDirectory string, fileExtension string) Store {
	return fileSystemStore{
		baseDirectory: baseDirectory,
		fileExtension: fileExtension,
	}
}

type fileSystemStore struct {
	baseDirectory string
	fileExtension string
}

func (s fileSystemStore) List(itemType string) ([]string, error) {
	if err := s.ensure(itemType); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(filepath.Join(s.baseDirectory, itemType))
	if err != nil {
		return []string{}, err
	}

	return names(s.storageFiles(files)), nil
}

func (s fileSystemStore) Save(itemType string, name string, data []byte) error {
	filename, err := s.fullyQualifiedName(itemType, name)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}

func (s fileSystemStore) Read(itemType string, name string) ([]byte, error) {
	filename, err := s.fullyQualifiedName(itemType, name)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, ErrRecordDoesNotExist
	}

	return ioutil.ReadFile(filename)
}

func (s fileSystemStore) Delete(itemType string, name string) error {
	filename, err := s.fullyQualifiedName(itemType, name)
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

func (s fileSystemStore) fileNameOf(itemType string, name string) string {
	return filepath.Join(s.baseDirectory, itemType, fmt.Sprintf("%s.%s", name, s.fileExtension))
}

func (s fileSystemStore) fullyQualifiedName(itemType string, name string) (string, error) {
	if err := s.ensure(itemType); err != nil {
		return "", err
	}
	return s.fileNameOf(itemType, name), nil
}

func (s fileSystemStore) ensure(itemType string) error {
	target := filepath.Join(s.baseDirectory, itemType)
	fi, err := os.Stat(target)
	if err == nil {
		if fi.IsDir() {
			return nil
		}
		return fmt.Errorf("storage path %s exists, but is not a directory", target)
	}
	return os.MkdirAll(target, os.ModePerm)
}

func (s fileSystemStore) storageFiles(files []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	ext := "." + s.fileExtension
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ext {
			result = append(result, file)
		}
	}
	return result
}

func names(files []os.FileInfo) []string {
	result := make([]string, 0)
	for _, file := range files {
		result = append(result, name(file.Name()))
	}
	return result
}

func name(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
