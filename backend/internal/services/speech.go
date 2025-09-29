package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type SpeechService interface {
	TextToSpeech(prompt string, voice VoiceEnum, speed float64) (*File, error)
	DownloadAudio(url string) ([]byte, error)
	SaveAudio(data []byte, filename string) error
}

type speechService struct {
	client      *FalClient
	maxAttempts int
	interval    time.Duration
}

func NewSpeechService(falClient *FalClient, maxAttempts int, interval time.Duration) SpeechService {
	return &speechService{
		client:      falClient,
		maxAttempts: maxAttempts,
		interval:    interval,
	}
}

// TextToSpeech converts text to speech using the fal.ai API
func (s *speechService) TextToSpeech(prompt string, voice VoiceEnum, speed float64) (*File, error) {
	if voice == "" {
		voice = VoiceAfHeart
	}
	if speed == 0 {
		speed = 1.0
	}

	req := EnglishRequest{
		Prompt: prompt,
		Voice:  voice,
		Speed:  speed,
	}

	result, err := s.client.submitAndWait(req, s.maxAttempts, s.interval)
	if err != nil {
		return nil, fmt.Errorf("failed to generate speech: %w", err)
	}

	return &result.Audio, nil
}

// DownloadAudio downloads audio data from a URL
func (s *speechService) DownloadAudio(url string) ([]byte, error) {
	return s.client.downloadAudioFile(url)
}

// SaveAudio saves audio data to a file
func (s *speechService) SaveAudio(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}

// VoiceEnum represents available voice options for the API
type VoiceEnum string

const (
	VoiceAfHeart   VoiceEnum = "af_heart"
	VoiceAfAlloy   VoiceEnum = "af_alloy"
	VoiceAfAoede   VoiceEnum = "af_aoede"
	VoiceAfBella   VoiceEnum = "af_bella"
	VoiceAfJessica VoiceEnum = "af_jessica"
	VoiceAfKore    VoiceEnum = "af_kore"
	VoiceAfNicole  VoiceEnum = "af_nicole"
	VoiceAfNova    VoiceEnum = "af_nova"
	VoiceAfRiver   VoiceEnum = "af_river"
	VoiceAfSarah   VoiceEnum = "af_sarah"
	VoiceAfSky     VoiceEnum = "af_sky"
	VoiceAmAdam    VoiceEnum = "am_adam"
	VoiceAmEcho    VoiceEnum = "am_echo"
	VoiceAmEric    VoiceEnum = "am_adam"
	VoiceAmFenrir  VoiceEnum = "am_fenrir"
	VoiceAmLiam    VoiceEnum = "am_liam"
	VoiceAmMichael VoiceEnum = "am_michael"
	VoiceAmOnyx    VoiceEnum = "am_onyx"
	VoiceAmPuck    VoiceEnum = "am_puck"
	VoiceAmSanta   VoiceEnum = "am_santa"
)

// File represents a file response from the API
type File struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	FileData    string `json:"file_data"`
}

// EnglishRequest represents a request for English text-to-speech
type EnglishRequest struct {
	Prompt string    `json:"prompt"`
	Voice  VoiceEnum `json:"voice"`
	Speed  float64   `json:"speed"`
}

// RequestResponse represents the initial response when submitting a request
type RequestResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// StatusResponse represents the status of a request
type StatusResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Progress  int    `json:"progress"`
}

// ResultResponse represents the final result of a request
type ResultResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Audio     File   `json:"audio"`
}

// FalClient represents a client for interacting with the fal.ai API
type FalClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewFalClient creates a new FalClient instance
func NewFalClient(apiKey string) *FalClient {
	return &FalClient{
		apiKey:  apiKey,
		baseURL: "https://queue.fal.run",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIError represents an error response from the API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// makeRequest makes an HTTP request to the fal.ai API
func (c *FalClient) makeRequest(method, url string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Key "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return resp, nil
}

// submitEnglishRequest submits a text-to-speech request for English text
func (c *FalClient) submitEnglishRequest(req EnglishRequest) (*RequestResponse, error) {
	url := fmt.Sprintf("%s/fal-ai/kokoro/american-english", c.baseURL)
	resp, err := c.makeRequest("POST", url, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// getRequestStatus checks the status of a request
func (c *FalClient) getRequestStatus(requestID string) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/fal-ai/kokoro/requests/%s/status", c.baseURL, requestID)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &result, nil
}

// getRequestResult retrieves the final result of a completed request
func (c *FalClient) getRequestResult(requestID string) (*ResultResponse, error) {
	url := fmt.Sprintf("%s/fal-ai/kokoro/requests/%s", c.baseURL, requestID)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode result response: %w", err)
	}

	return &result, nil
}

// waitForResult waits for a request to complete and returns the result
func (c *FalClient) waitForResult(requestID string, maxAttempts int, interval time.Duration) (*ResultResponse, error) {
	for i := 0; i < maxAttempts; i++ {
		status, err := c.getRequestStatus(requestID)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Status check %d: %s", i+1, status.Status)
		if status.Progress > 0 {
			fmt.Printf(" (%d%%)", status.Progress)
		}
		fmt.Println()

		switch status.Status {
		case "COMPLETED":
			return c.getRequestResult(requestID)
		case "FAILED":
			return nil, fmt.Errorf("request failed: %s", status.Message)
		case "IN_PROGRESS", "QUEUED", "IN_QUEUE":
			time.Sleep(interval)
			continue
		default:
			return nil, fmt.Errorf("unknown status: %s", status.Status)
		}
	}

	return nil, fmt.Errorf("request did not complete within %d attempts", maxAttempts)
}

// submitAndWait submits a request and waits for it to complete
func (c *FalClient) submitAndWait(req EnglishRequest, maxAttempts int, interval time.Duration) (*ResultResponse, error) {
	submitResp, err := c.submitEnglishRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to submit request: %w", err)
	}

	return c.waitForResult(submitResp.RequestID, maxAttempts, interval)
}

// WithRetry adds retry logic to any operation
func WithRetry(operation func() error, maxAttempts int, backoff time.Duration) error {
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't wait on the last attempt
		if i < maxAttempts-1 {
			time.Sleep(backoff * time.Duration(i+1))
		}
	}

	return fmt.Errorf("operation failed after %d attempts, last error: %w", maxAttempts, lastErr)
}

// downloadAudioFile downloads an audio file from the given URL
func (c *FalClient) downloadAudioFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download audio file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download audio file, status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// CreateEnglishRequest creates a new English request with default values
func CreateEnglishRequest(prompt string, voice VoiceEnum, speed float64) EnglishRequest {
	if voice == "" {
		voice = VoiceAfHeart
	}
	if speed == 0 {
		speed = 1.0
	}
	return EnglishRequest{
		Prompt: prompt,
		Voice:  voice,
		Speed:  speed,
	}
}

// StitchAudioFiles combines multiple audio files into a single file using ffmpeg
func StitchAudioFiles(inputFiles []string, outputFile string) error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// Create a temporary file list for ffmpeg
	listFile, err := os.CreateTemp("", "audio_list_*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(listFile.Name())
	defer listFile.Close()

	// Write file list in ffmpeg format
	for _, file := range inputFiles {
		if _, err := listFile.WriteString(fmt.Sprintf("file '%s'\n", file)); err != nil {
			return fmt.Errorf("failed to write to list file: %w", err)
		}
	}

	// Build ffmpeg command
	cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", listFile.Name(), "-c", "copy", outputFile)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
