package handlers

import (
	"net/http"
	"io"
	"encoding/json"
	"fmt"
	"crypto/md5"
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
		io.WriteString(w, "Method GET")
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
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	hashDigest := getHash(s.PlainText)
	secretResponse := SecretResponse{Id: hashDigest}

	res, err := json.Marshal(&secretResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}