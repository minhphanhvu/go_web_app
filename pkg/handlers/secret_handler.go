package handlers

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"crypto/md5"
	"strings"

	"github.com/minhphanhvu/go_web_app/pkg/filestore"
	"github.com/minhphanhvu/go_web_app/pkg/types"
)

type Secret struct {
	PlainText string `json:"plain_text"`
}

type SecretResponse struct {
	Id string `json:"id"`
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "POST" {
		createSecret(w, r)
	} else if method == "GET" {
		getSecret(w, r)
	} else {
		w.WriteHeader(405);
		io.WriteString(w, "Method is not allowed.")
	}
}

func getHash(plainText string) string {
	text := []byte(plainText)
	return fmt.Sprintf("%x", md5.Sum(text))
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	s := Secret{}
	err := json.NewDecoder(r.Body).Decode(&s) // Unmarshal the request body into s

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if len(s.PlainText) == 0 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	hashDigest := getHash(s.PlainText)
	secretResponse := SecretResponse{Id: hashDigest}

	// Write secret to the file
	secretData := types.SecretData{Id: hashDigest, Secret: s.PlainText}
	// Write method read all data from data.json, rewrite it with 
	// existing and new data
	err = filestore.FileStoreConfig.Fs.Write(secretData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal secretResponse to send it back to the client
	res, err := json.Marshal(&secretResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// Retrieve the value from specified id -> delete that from data.json
func getSecret(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	id = strings.TrimPrefix(id, "/")

	if len(id) == 0 {
		http.Error(w, "No Secret ID specified", http.StatusBadRequest)
		return
	}

	data, err := filestore.FileStoreConfig.Fs.Read(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.SecretResponse{}
	response.Data = data

	jRes, err := json.Marshal(&response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(response.Data) == 0 {
		w.WriteHeader(404)
	}
	w.Write(jRes)
}