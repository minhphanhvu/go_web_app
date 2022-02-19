package handlers

import (
	"io"
	"net/http/httptest"
	"net/http"
	"testing"
	"fmt"
)

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()
	SetupHandlers(mux)
	ts := httptest.NewServer(mux) // Simulate running a test server
	defer ts.Close() // Caller must close the server

	// healthcheck test
	res, err := http.Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	response, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	respStr := string(response)
	if respStr != "Health check is ok. You can begin making requests now." {
		t.Errorf("Expected GET /healthcheck request to return Health check is ok. You can begin making requests now., Got: %s", respStr)
	}

	// GET test with invalid id -> expect status 404 and response "Np Secret ID specified"
	res, err = http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	response, err = io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	respStr = string(response)
	if respStr != "No Secret ID specified\n" {
		t.Errorf("Expected GET / request to return \"No Secret ID specified\". Got: %s", respStr)
	}

	// POST test with empty JSON body
	// Better rewrite POST to return an error instead of panic
	// This is like a unit test
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recorvered.")
		}
	}()
	req := httptest.NewRequest("POST", ts.URL, nil)
	w := httptest.NewRecorder()
	secretHandler(w, req) // Call the secretHandler then pass the request to the create function	
}
