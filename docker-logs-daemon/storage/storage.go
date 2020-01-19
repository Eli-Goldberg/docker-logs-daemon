package storage

import "io"

// Storage implements Storage reading and writing
type Storage interface {
	Writer(streamName string) (io.WriteCloser, error)
	Reader(streamName string) (io.ReadCloser, error)
	List() ([]string, error)
}
