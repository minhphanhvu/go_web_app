package main

import (
	"os"
	"net/http"
	"github.com/minhphanhvu/go_web_app/pkg/filestore"
	// "github.com/minhphanhvu/go_web_app/pkg/handlers"
	"log"
)

func main() {
	mux := http.NewServeMux()
	// handlers.SetupHandlers(mux)

	DATA_FILE_PATH := os.Getenv("DATA_FILE_PATH")
	if len(DATA_FILE_PATH) == 0 {
		log.Fatal("Please specify the PATH")
		os.Exit(1)
	}
	filestore.Init(DATA_FILE_PATH)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server cannot start on port :8080, Error: %v", err)
	}
}