package storage

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// FileStorage implementation
type FileStorage struct {
	path string
}

// FileStorageType is used to uniquely identify the "File" storage type
const FileStorageType = "File"

// DefaultFileStoragePath is the default file storage path
const DefaultFileStoragePath = "./logs"

// NewFileStorage Creates a new FileStorage
func NewFileStorage(path string) (*FileStorage, error) {

	absPath, pathErr := filepath.Abs(path)
	if pathErr != nil {
		return nil, pathErr
	}

	if dirErr := ensureDir(absPath); dirErr != nil {
		return nil, dirErr
	}

	return &FileStorage{
		path: absPath,
	}, nil
}

// Writer returns an io.WriteCloser stream to pass logs into it
// It's up to the caller to close the stream.
func (fs FileStorage) Writer(streamName string) (io.WriteCloser, error) {
	streamPath := fs.path + "/" + streamName

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(streamPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return f, err
}

// Reader returns an io.ReadCloser stream to read logs from
func (fs FileStorage) Reader(streamName string) (io.ReadCloser, error) {
	streamPath := fs.path + "/" + streamName

	f, err := os.Open(streamPath)
	return f, err
}

// List returns all possible stream names
func (fs FileStorage) List() ([]string, error) {
	var fileNames []string

	files, err := ioutil.ReadDir(fs.path)
	if err != nil {
		log.Fatal(err)
		return fileNames, err
	}

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModePerm)

	if err == nil || os.IsExist(err) {
		return nil
	}

	return err
}
