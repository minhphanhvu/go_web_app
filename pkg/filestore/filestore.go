package filestore

import (
	"os"
	"sync"
	"io"
	"log"
	"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/minhphanhvu/go_web_app/pkg/types"
	"golang.org/x/crypto/scrypt"
)

type fileStore struct {
	Mu sync.Mutex
	Store map[string][]byte
}

var FileStoreConfig struct {
	DataFilePath string
	GCM 				 cipher.AEAD
	Nonce 			 []byte
	Fs           fileStore
}

func Init(dataFilePath string, password, salt string) error {
	_, err := os.Stat(dataFilePath) // return information about the file, error if path does not exist.

	if err != nil {
		_, err := os.Create(dataFilePath) // if file not exists, create it.
		if err != nil {
			return err
		}
	}

	gcm, nonce, err := initCrypto(password, salt)
	if err != nil {
		return err
	}
	FileStoreConfig.GCM = gcm
	FileStoreConfig.Nonce = nonce

	FileStoreConfig.Fs = fileStore{Mu: sync.Mutex{}, Store: make(map[string][]byte)}
	FileStoreConfig.DataFilePath = dataFilePath
	return nil
}

// For encryption/descryption
// gcm, nonce returned are used to encrypt and decrypt data
func initCrypto(password, salt string) (cipher.AEAD, []byte, error) {
	key, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	return gcm, nonce, err
}

func encrypt(data string) []byte {

	if _, err := io.ReadFull(rand.Reader, FileStoreConfig.Nonce); err != nil {
		log.Fatal(err)
	}
	encryptedData := FileStoreConfig.GCM.Seal(FileStoreConfig.Nonce, FileStoreConfig.Nonce, []byte(data), nil)
	return encryptedData
}

func decrypt(encData []byte) ([]byte, error) {
	nonce := encData[:FileStoreConfig.GCM.NonceSize()]
	encData = encData[FileStoreConfig.GCM.NonceSize():]
	data, err := FileStoreConfig.GCM.Open(nil, nonce, encData, nil)
	if err != nil {
					return nil, err
	}
	return data, err
}

func (fileStore *fileStore) Write(data types.SecretData) error {
	fileStore.Mu.Lock()
	defer fileStore.Mu.Unlock()

	err := fileStore.ReadFromFile()
	if err != nil {
		return err
	}

	fileStore.Store[data.Id] = encrypt(data.Secret)
	return fileStore.WriteToFile()
}

func (fileStore *fileStore) Read(id string) (string, error) {
	fileStore.Mu.Lock()
	defer fileStore.Mu.Unlock()

	err := fileStore.ReadFromFile()
	if err != nil {
		return "", err
	}

	encData, ok := fileStore.Store[id]
	if !ok {
		return "", nil
	}

	data, err := decrypt(encData)
	if err != nil {
		return "", err
	}
	delete(fileStore.Store, id)
	fileStore.WriteToFile()

	return string(data), nil
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
