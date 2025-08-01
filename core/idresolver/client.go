package idresolver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// IDMSUClient implements StudentIDResolver by calling the idmsu service
type IDMSUClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewIDMSUClient creates a new client for the idmsu service
func NewIDMSUClient() *IDMSUClient {
	baseURL := os.Getenv("IDMSU_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	return &IDMSUClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// ResolveBatch implements StudentIDResolver interface by calling the idmsu service
func (c *IDMSUClient) ResolveBatch(ctx context.Context, req []ResolveRequestItem) ([]ResolveResponseItem, error) {
	const batchSize = 1000
	var allResults []ResolveResponseItem

	for i := 0; i < len(req); i += batchSize {
		end := i + batchSize
		if end > len(req) {
			end = len(req)
		}
		batch := req[i:end]

		var err error
		var reqBody []byte
		var httpReq *http.Request
		var resp *http.Response
		var respBody []byte

		reqBody, err = json.Marshal(map[string][]ResolveRequestItem{"items": batch})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		url := fmt.Sprintf("%s/api/v1/resolve", c.baseURL)

		const maxRetries = 10
		// Retry indefinitely until the service responds with 200 OK or the context is cancelled.
		for attempt := 0; ; attempt++ {
			httpReq, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
			if err != nil {
				return nil, fmt.Errorf("failed to create HTTP request: %w", err)
			}

			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Accept", "application/json")

			resp, err = c.httpClient.Do(httpReq)
			if err != nil {
				return nil, fmt.Errorf("failed to make HTTP request: %w", err)
			}

			if resp.StatusCode == http.StatusOK {
				break
			}
			if resp.StatusCode == http.StatusServiceUnavailable {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				delay := 1 * time.Minute
				time.Sleep(delay)
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("idmsu service returned status %d: %s", resp.StatusCode, string(body))
		}

		// resp.StatusCode is guaranteed to be 200 OK here; no further readiness check required.
		defer resp.Body.Close()

		// Parse response
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var batchResult []ResolveResponseItem
		if err = json.Unmarshal(respBody, &batchResult); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allResults = append(allResults, batchResult...)
	}

	return allResults, nil
}

// Health checks if the idmsu service is healthy
func (c *IDMSUClient) Health(ctx context.Context) error {
	var err error
	var req *http.Request
	var resp *http.Response

	url := fmt.Sprintf("%s/api/v1/health", c.baseURL)
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make health check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("idmsu service health check failed with status %d", resp.StatusCode)
	}

	return nil
}