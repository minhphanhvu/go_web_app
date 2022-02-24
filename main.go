package main

import (
	"os"
	"net/http"
	"github.com/minhphanhvu/go_web_app/pkg/filestore"
	"github.com/minhphanhvu/go_web_app/pkg/handlers"
	"log"
)

func main() {
	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	DATA_FILE_PATH := os.Getenv("DATA_FILE_PATH")
	if len(DATA_FILE_PATH) == 0 {
		log.Fatal("Please specify the PATH")
		os.Exit(1)
	}
	PASSWORD := os.Getenv("PASSWORD")
	if len(PASSWORD) == 0 {
		log.Fatal("Please specify the PASSWORD for encryption/decryption")
		os.Exit(1)
	}
	SALT := os.Getenv("SALT")
	if len(SALT) == 0 {
		log.Fatal("Please specify the SALT for encryption/decryption")
		os.Exit(1)
	}
	filestore.Init(DATA_FILE_PATH, PASSWORD, SALT)

	err := http.ListenAndServe(":8080", mux)
	// if server not successfully, terminal will log message
	if err != nil {
		log.Fatalf("Server cannot start on port :8080, Error: %v", err)
	}
}