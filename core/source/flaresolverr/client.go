package flaresolverr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Client represents a FlareSolverr client wrapper
type Client struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

// Request represents a FlareSolverr API request
type Request struct {
	Cmd        string                 `json:"cmd"`
	URL        string                 `json:"url,omitempty"`
	MaxTimeout int                    `json:"maxTimeout,omitempty"`
	Session    string                 `json:"session,omitempty"`
	PostData   map[string]interface{} `json:"postData,omitempty"`
	Headers    map[string]string      `json:"headers,omitempty"`
}

// Response represents a FlareSolverr API response
type Response struct {
	Status   string    `json:"status"`
	Message  string    `json:"message"`
	Session  string    `json:"session,omitempty"`
	Solution *Solution `json:"solution,omitempty"`
}

// Solution contains the actual HTTP response data
type Solution struct {
	URL      string                   `json:"url"`
	Status   int                      `json:"status"`
	Headers  map[string]string        `json:"headers"`
	Response string                   `json:"response"`
	Cookies  []map[string]interface{} `json:"cookies"`
}

// GetResponse represents a simplified response for GET requests
type GetResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
	URL        string
	Error      error
}

var (
	clientInstance *Client
	clientOnce     sync.Once
	clientError    error
)

// GetClient returns a singleton FlareSolverr client instance
func GetClient() (*Client, error) {
	clientOnce.Do(func() {
		clientInstance, clientError = initializeClient()
	})
	return clientInstance, clientError
}

// initializeClient creates and configures a new FlareSolverr client
func initializeClient() (*Client, error) {
	// Get configuration from environment variables
	baseURL := os.Getenv("FLARESOLVERR_URL")
	if baseURL == "" {
		// Default to localhost for development, but allow override
		baseURL = "http://localhost:8191"
	}

	timeoutStr := os.Getenv("FLARESOLVERR_TIMEOUT_MS")
	timeoutMs := 60000 // Default timeout
	if timeoutStr != "" {
		if parsed, err := strconv.Atoi(timeoutStr); err == nil {
			timeoutMs = parsed
		}
	}

	timeout := time.Duration(timeoutMs) * time.Millisecond

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: timeout + (10 * time.Second), // Add buffer for FlareSolverr overhead
	}

	client := &Client{
		baseURL:    baseURL,
		timeout:    timeout,
		httpClient: httpClient,
	}

	// Test connectivity
	if err := client.testConnection(); err != nil {
		return nil, fmt.Errorf("failed to connect to FlareSolverr at %s: %w", baseURL, err)
	}

	return client, nil
}

// testConnection verifies that FlareSolverr is accessible
func (c *Client) testConnection() error {
	resp, err := c.httpClient.Get(c.baseURL + "/v1")
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMethodNotAllowed {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Get performs a GET request through FlareSolverr
func (c *Client) Get(url string) (*GetResponse, error) {
	return c.GetWithHeaders(url, nil)
}

// GetWithHeaders performs a GET request with custom headers through FlareSolverr
func (c *Client) GetWithHeaders(url string, headers map[string]string) (*GetResponse, error) {
	request := Request{
		Cmd:        "request.get",
		URL:        url,
		MaxTimeout: int(c.timeout.Milliseconds()),
		Headers:    headers,
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return &GetResponse{Error: err}, err
	}

	if response.Status != "ok" {
		err := fmt.Errorf("FlareSolverr error: %s", response.Message)
		return &GetResponse{Error: err}, err
	}

	if response.Solution == nil {
		err := fmt.Errorf("no solution in FlareSolverr response")
		return &GetResponse{Error: err}, err
	}

	return &GetResponse{
		StatusCode: response.Solution.Status,
		Body:       response.Solution.Response,
		Headers:    response.Solution.Headers,
		URL:        response.Solution.URL,
	}, nil
}

// GetWithSession performs a GET request using a specific session
func (c *Client) GetWithSession(url string, sessionID string, headers map[string]string) (*GetResponse, error) {
	request := Request{
		Cmd:        "request.get",
		URL:        url,
		MaxTimeout: int(c.timeout.Milliseconds()),
		Session:    sessionID,
		Headers:    headers,
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return &GetResponse{Error: err}, err
	}

	if response.Status != "ok" {
		err := fmt.Errorf("FlareSolverr error: %s", response.Message)
		return &GetResponse{Error: err}, err
	}

	if response.Solution == nil {
		err := fmt.Errorf("no solution in FlareSolverr response")
		return &GetResponse{Error: err}, err
	}

	return &GetResponse{
		StatusCode: response.Solution.Status,
		Body:       response.Solution.Response,
		Headers:    response.Solution.Headers,
		URL:        response.Solution.URL,
	}, nil
}

// CreateSession creates a new FlareSolverr session
func (c *Client) CreateSession() (string, error) {
	request := Request{
		Cmd: "sessions.create",
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	if response.Status != "ok" {
		return "", fmt.Errorf("session creation failed: %s", response.Message)
	}

	return response.Session, nil
}

// DestroySession destroys a FlareSolverr session
func (c *Client) DestroySession(sessionID string) error {
	request := Request{
		Cmd:     "sessions.destroy",
		Session: sessionID,
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("session destruction failed: %s", response.Message)
	}

	return nil
}

// sendRequest sends a request to the FlareSolverr API
func (c *Client) sendRequest(req Request) (*Response, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// IsAvailable checks if FlareSolverr service is available
func (c *Client) IsAvailable() bool {
	return c.testConnection() == nil
}

// GetBaseURL returns the configured FlareSolverr base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// GetTimeout returns the configured timeout duration
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}

// SafeGet performs a GET request with graceful fallback handling
// Returns an error that can be checked for FlareSolverr availability
func SafeGet(url string) (*GetResponse, error) {
	client, err := GetClient()
	if err != nil {
		return &GetResponse{Error: err}, fmt.Errorf("FlareSolverr unavailable: %w", err)
	}

	return client.Get(url)
}

// SafeGetWithHeaders performs a GET request with headers and graceful fallback handling
func SafeGetWithHeaders(url string, headers map[string]string) (*GetResponse, error) {
	client, err := GetClient()
	if err != nil {
		return &GetResponse{Error: err}, fmt.Errorf("FlareSolverr unavailable: %w", err)
	}

	return client.GetWithHeaders(url, headers)
}

// IsFlareSolverrError checks if an error is related to FlareSolverr unavailability
func IsFlareSolverrError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return contains(errStr, "FlareSolverr unavailable") ||
		contains(errStr, "connection test failed") ||
		contains(errStr, "failed to connect to FlareSolverr")
}

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) >= 0)))
}

// indexOfSubstring finds the index of a substring in a string
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
