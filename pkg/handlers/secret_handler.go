package handlers

import (
	"net/http"
	"io"
)

func secretHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		io.WriteString(w, "Method Post")
	} else if r.Method == "GET" {
		io.WriteString(w, "Method GET")
	} else {
		w.WriteHeader(405);
		io.WriteString(w, "Method is not allowed.")
	}
}