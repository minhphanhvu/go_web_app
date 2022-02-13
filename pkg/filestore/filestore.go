package filestore

import (
	"os"
	"sync"
	"github.com/minhphanhvu/go_web_app/pkg/types"
	"io"
	"log"
	"encoding/json"
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

	FileStoreConfig.Fs = fileStore{Mu: sync.Mutex{}, Store: make(map[string]string)}
	FileStoreConfig.DataFilePath = dataFilePath
	return nil
}

func (fileStore *fileStore) Write(data types.SecretData) error {
	fileStore.Mu.Lock()
	defer fileStore.Mu.Unlock()

	err := fileStore.ReadFromFile()
	if err != nil {
		return err
	}

	fileStore.Store[data.Id] = data.Secret
	return fileStore.WriteToFile()
}

func (fileStore *fileStore) ReadFromFile() error {
	file, err := os.Open(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(jsonData) != 0 {
		json.Unmarshal(jsonData, &fileStore.Store)
	}
	return nil
}

func (fileStore *fileStore) WriteToFile() error {
	var file *os.File
	jsonData, err := json.Marshal(fileStore.Store)
	if err != nil {
		return err
	}
	file, err = os.Create(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	return err
}