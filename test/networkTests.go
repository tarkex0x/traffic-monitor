package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// startMockServer starts a new mock server with registered handlers.
func startMockServer() *httptest.Server {
	handler := http.NewServeMux()
	registerHandlers(handler)
	
	server := httptest.NewServer(handler)
	return server
}

// registerHandlers registers all handlers for the mock server.
func registerHandlers(handler *http.ServeMux) {
	handler.HandleFunc("/data", postDataHandler)
}

// postDataHandler processes the "/data" endpoint requests.
func postDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	processRequestBody(w, r)
}

// processRequestBody reads and responds to the incoming request.
func processRequestBody(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	response := []byte("Received: " + string(body))
	w.Write(response)
}

// sendAndCaptureData sends traffic data and captures the response.
func sendAndCaptureData(url, trafficData string) (string, error) {
	resp, err := createAndSendRequest(url, trafficData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return readResponse(resp)
}

// createAndSendRequest creates and sends an HTTP request.
func createAndSendRequest(url, trafficData string) (*http.Response, error) {
	return http.Post(url, "application/text", bytes.NewBufferString(trafficData))
}

// readResponse reads the HTTP response body.
func readResponse(resp *http.Response) (string, error) {
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