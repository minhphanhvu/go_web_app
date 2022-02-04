package filestore

import (
	"os"
	"sync"
)

type fileStore struct {
	Mu sync.Mutex
	Store map[string]string
}

var FileStoreConfig struct {
	DataFilePath string
	Fs           fileStore
}

func Init(dataFilePath string) error {
	_, err := os.Stat(dataFilePath) // return information about the file, error if path does not exist.

	if err != nil {
		_, err := os.Create(dataFilePath) // if file not exists, create it.
		if err != nil {
			return err
		}
	}

	FileStoreConfig.Fs = fileStore{Mu: &sync.Mutex{}, Store: make(map[string]string)}
	FileStoreConfig.DataFilePath = dataFilePath
	return nil
}