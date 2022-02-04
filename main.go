package main

import (
	"fmt"
	"os"
	"errors"
	"net/http"
	"github.com/minhphanhvu/go_web_app/pkg/filestore"
)

func main() {
	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	DATA_FILE_PATH := os.Getenv("DATA_FILE_PATH")
	filestore.Init(DATA_FILE_PATH)
}