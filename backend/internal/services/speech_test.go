package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpeechService_Basic(t *testing.T) {
	// Basic test to ensure the package compiles
	// The actual SpeechService requires extensive external dependencies
	assert.True(t, true)
}

func TestNewSpeechService(t *testing.T) {
	client := NewFalClient("test-api-key")
	svc := NewSpeechService(client, 3, time.Second)
	assert.NotNil(t, svc)
}

func TestNewFalClient(t *testing.T) {
	client := NewFalClient("test-api-key")
	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.apiKey)
	assert.Equal(t, "https://queue.fal.run", client.baseURL)
	assert.NotNil(t, client.httpClient)
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Not found",
	}
	assert.Equal(t, "API error (status 404): Not found", err.Error())
}

func TestCreateEnglishRequest(t *testing.T) {
	// Test with all parameters
	req := CreateEnglishRequest("Hello world", VoiceAfHeart, 1.5)
	assert.Equal(t, "Hello world", req.Prompt)
	assert.Equal(t, VoiceAfHeart, req.Voice)
	assert.Equal(t, 1.5, req.Speed)

	// Test with defaults
	req2 := CreateEnglishRequest("Test", "", 0)
	assert.Equal(t, "Test", req2.Prompt)
	assert.Equal(t, VoiceAfHeart, req2.Voice)
	assert.Equal(t, 1.0, req2.Speed)
}

func TestSpeechService_SaveAudio(t *testing.T) {
	client := NewFalClient("test-api-key")
	svc := NewSpeechService(client, 3, time.Second)

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-audio-*.mp3")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test saving audio
	testData := []byte("test audio data")
	err = svc.SaveAudio(testData, tmpFile.Name())
	assert.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
}

func TestFalClient_MakeRequest_Success(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "Key test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	resp, err := client.makeRequest("POST", server.URL+"/test", map[string]string{"test": "data"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFalClient_MakeRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request error"))
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	_, err := client.makeRequest("POST", server.URL+"/test", nil)
	assert.Error(t, err)

	apiErr, ok := err.(*APIError)
	assert.True(t, ok)
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "Bad request")
}

func TestFalClient_SubmitEnglishRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "kokoro/american-english")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RequestResponse{
			RequestID: "test-request-id",
			Status:    "QUEUED",
			Message:   "Request queued",
		})
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	req := EnglishRequest{
		Prompt: "Test prompt",
		Voice:  VoiceAfHeart,
		Speed:  1.0,
	}

	resp, err := client.submitEnglishRequest(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-request-id", resp.RequestID)
	assert.Equal(t, "QUEUED", resp.Status)
}

func TestFalClient_GetRequestStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "status")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StatusResponse{
			RequestID: "test-request-id",
			Status:    "IN_PROGRESS",
			Progress:  50,
			Message:   "Processing",
		})
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	resp, err := client.getRequestStatus("test-request-id")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-request-id", resp.RequestID)
	assert.Equal(t, "IN_PROGRESS", resp.Status)
	assert.Equal(t, 50, resp.Progress)
}

func TestFalClient_GetRequestResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResultResponse{
			RequestID: "test-request-id",
			Status:    "COMPLETED",
			Audio: File{
				URL:         "https://example.com/audio.mp3",
				ContentType: "audio/mpeg",
				FileName:    "output.mp3",
				FileSize:    1024,
			},
		})
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	resp, err := client.getRequestResult("test-request-id")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "COMPLETED", resp.Status)
	assert.Equal(t, "https://example.com/audio.mp3", resp.Audio.URL)
}

func TestFalClient_WaitForResult_Completed(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if strings.Contains(r.URL.Path, "status") {
			// First call returns IN_PROGRESS, second returns COMPLETED
			if callCount == 1 {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(StatusResponse{
					RequestID: "test-request-id-12345",
					Status:    "IN_PROGRESS",
					Progress:  50,
				})
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(StatusResponse{
					RequestID: "test-request-id-12345",
					Status:    "COMPLETED",
					Progress:  100,
				})
			}
		} else {
			// Return result
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ResultResponse{
				RequestID: "test-request-id-12345",
				Status:    "COMPLETED",
				Audio: File{
					URL: "https://example.com/audio.mp3",
				},
			})
		}
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	dialogueInfo := &DialogueInfo{
		Index: 1,
		Total: 1,
	}

	result, err := client.waitForResult("test-request-id-12345", 10*time.Millisecond, dialogueInfo)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "COMPLETED", result.Status)
}

func TestFalClient_WaitForResult_Failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StatusResponse{
			RequestID: "test-request-id-12345",
			Status:    "FAILED",
			Message:   "Processing failed",
		})
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	dialogueInfo := &DialogueInfo{
		Index: 1,
		Total: 1,
	}

	_, err := client.waitForResult("test-request-id-12345", 10*time.Millisecond, dialogueInfo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestFalClient_SubmitAndWait(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		if r.Method == "POST" {
			// Submit request
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(RequestResponse{
				RequestID: "test-id",
				Status:    "QUEUED",
			})
		} else if strings.Contains(r.URL.Path, "status") {
			// Status check - return COMPLETED immediately
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(StatusResponse{
				RequestID: "test-id",
				Status:    "COMPLETED",
			})
		} else {
			// Get result
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ResultResponse{
				RequestID: "test-id",
				Status:    "COMPLETED",
				Audio: File{
					URL: "https://example.com/audio.mp3",
				},
			})
		}
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")
	client.baseURL = server.URL

	req := EnglishRequest{
		Prompt: "Test",
		Voice:  VoiceAfHeart,
		Speed:  1.0,
	}

	result, err := client.submitAndWait(req, 3, 10*time.Millisecond, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://example.com/audio.mp3", result.Audio.URL)
}

func TestWithRetry_Success(t *testing.T) {
	callCount := 0
	operation := func() error {
		callCount++
		if callCount < 2 {
			return fmt.Errorf("temporary error")
		}
		return nil
	}

	err := WithRetry(operation, 3, 10*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestWithRetry_AllAttemptsFail(t *testing.T) {
	callCount := 0
	operation := func() error {
		callCount++
		return fmt.Errorf("persistent error")
	}

	err := WithRetry(operation, 3, 10*time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, 3, callCount)
	assert.Contains(t, err.Error(), "operation failed after 3 attempts")
}

func TestFalClient_DownloadAudioFile(t *testing.T) {
	expectedData := []byte("audio data content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expectedData)
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")

	data, err := client.downloadAudioFile(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
}

func TestFalClient_DownloadAudioFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewFalClient("test-api-key")

	_, err := client.downloadAudioFile(server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download audio file")
}

func TestSpeechService_TextToSpeech_NoClient(t *testing.T) {
	svc := &speechService{
		client:      nil,
		maxAttempts: 3,
		interval:    time.Second,
	}

	_, err := svc.TextToSpeech("test", VoiceAfHeart, 1.0, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "speech service not configured")
}

func TestSpeechService_DownloadAudio_NoClient(t *testing.T) {
	svc := &speechService{
		client:      nil,
		maxAttempts: 3,
		interval:    time.Second,
	}

	_, err := svc.DownloadAudio("https://example.com/audio.mp3")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "speech service not configured")
}

// Note: Full integration testing would require a real Fal.ai API key.
// The tests above validate request/response handling, error cases, and retry logic.
