package handlers

import (
	"testing"
	"fmt"
	"path/filepath"
	"net/http/httptest"
	"strings"
	"io"
	"os"

	"github.com/minhphanhvu/go_web_app/pkg/filestore"
)
func TestHandlers(t *testing.T) {
	// CreateSecret handler test
	testCases := []struct {
		requestBody 					 string
		expectedHTTPStatusCode int
		expectedBody 					 string
	}{
		{
			requestBody: 					 "{\"plaintext\": \"some secret\"}",
			expectedHTTPStatusCode: 400,
			expectedBody: 				 "Invalid Request\n",
		},
		{
			requestBody: 					 "{\"plain_text\": \"some secret\"}",
			expectedHTTPStatusCode: 200,
			expectedBody:          fmt.Sprintf("{\"id\":\"%s\"}", getHash("some secret")),
		},
	}

	filestore.Init(filepath.Join(t.TempDir(), "test_data.json"))

	for _, tc := range testCases {
		req := httptest.NewRequest("POST", "/", strings.NewReader(tc.requestBody))
		w := httptest.NewRecorder()
		secretHandler(w, req)

		response := w.Result()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != tc.expectedHTTPStatusCode {
			t.Errorf("Expected Response Status to be: %d, Got: %d", tc.expectedHTTPStatusCode, response.StatusCode)

		}
		if string(body) != tc.expectedBody {
			t.Errorf("Expected Response body to be: %s, Got: %s", tc.expectedBody, string(body))
		}
	}

	// GET request test

	// A premade secretID and secret inside a temp directory
	testSecretId := "7a819afa983d454b3a368c1422ba853c"
	testSecret := "My super secret1234151"
	testDataContents := fmt.Sprintf(`{"%s":"%s"}`, testSecretId, testSecret)
	testDataFilePath := filepath.Join(t.TempDir(), "test_data.json")
	f, err := os.Create(testDataFilePath)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(testDataContents)

	// Should be able to retrieve an existing id
	filestore.Init(testDataFilePath)
	expectedResponseBody := fmt.Sprintf(`{"data":"%s"}`, testSecret)

	req := httptest.NewRequest("GET", "/"+testSecretId, nil)
	w := httptest.NewRecorder()
	secretHandler(w, req)

	response := w.Result()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Expected Response Status to be: %d, Got: %d", 200, response.StatusCode)

	}
	if string(body) != expectedResponseBody {
		t.Errorf("Expected Response body to be: %s, Got: %s", expectedResponseBody, string(body))
	}
}