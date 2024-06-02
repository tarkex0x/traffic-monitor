package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func startTestServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			response := []byte("Received: " + string(body))
			w.Write(response)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})
	server := httptest.NewServer(handler)
	return server
}

func analyzeTrafficAndCaptureData(url, data string) (string, error) {
	resp, err := http.Post(url, "application/text", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func TestTrafficAnalysisAndDataCapturing(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	url := server.URL + "/data"

	mockData := "Hello, Network!"

	receivedData, err := analyzeTrafficAndCaptureData(url, mockData)
	if err != nil {
		t.Errorf("Failed to analyze traffic and capture data: %v", err)
	}

	expectedResponse := "Received: " + mockData
	if receivedData != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, receivedData)
	}
}

func setupFromEnv() {
	os.Setenv("NETWORK_ANALYSIS_ENDPOINT", "http://example.com/data")
}

func TestEnvBasedConfig(t *testinging.T) {
	setupFromParam()

	endpoint := os.Getenv("NETWORK_ANALYSIS_ENDPOINT")
	if endpoint == "" {
		t.Error("Environment variable NETWORK_ANALYSIS_ENDPOINT is not set")
	}
}