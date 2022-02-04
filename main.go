package main

import (
	"fmt"
	"os"
	"errors"
	"net/http"
	"github.com/minhphanhvu/go_web_app/pkg/filestore"
	"log"
)

func main() {
	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	DATA_FILE_PATH := os.Getenv("DATA_FILE_PATH")
	if len(DATA_FILE_PATH) == 0 {
		log.Fatal("Please specify the PATH")
	}
	filestore.Init(DATA_FILE_PATH)
}