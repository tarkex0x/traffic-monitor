package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func startMockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		response := []byte("Received: " + string(body))
		w.Write(response)
	})
	server := httptest.NewServer(handler)
	return server
}

func sendAndCaptureData(url, trafficData string) (string, error) {
	resp, err := http.Post(url, "application/text", bytes.NewBufferString(trafficData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func TestTrafficDataSendingAndCapturing(t *testing.T) {
	mockServer := startMockServer()
	defer mockServer.Close()

	mockServerURL := mockServer.URL + "/data"
	testTrafficData := "Hello, Network!"

	responseData, err := sendAndCaptureData(mockServerURL, testTrafficData)
	if err != nil {
		t.Errorf("Failed to send and capture data: %v", err)
	}

	expectedResponse := "Received: " + testTrafficData
	if responseData != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, responseData)
	}
}

func initializeFromEnvironmentVars() {
	os.Setenv("NETWORK_ANALYSIS_ENDPOINT", "http://example.com/data")
}

func TestEnvironmentVariableConfiguration(t *testing.T) {
	initializeFromEnvironmentVars()

	networkAnalysisEndpoint := os.Getenv("NETWORK_ANALYSIS_ENDPOINT")
	if networkAnalysisEndpoint == "" {
		t.Error("Environment variable NETWORK_ANALYSIS_ENDPOINT is not set")
	}
}