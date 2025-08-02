//go:build legacyclient
// +build legacyclient

package idresolver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
			Timeout: 15 * time.Minute,
		},
	}
}

// ResolveBatch implements StudentIDResolver interface by calling the idmsu service
func (c *IDMSUClient) ResolveBatch(ctx context.Context, req []ResolveRequestItem) ([]ResolveResponseItem, error) {
	var (
		err      error
		reqBody  []byte
		httpReq  *http.Request
		resp     *http.Response
		respBody []byte
	)

	reqBody, err = json.Marshal(map[string][]ResolveRequestItem{"items": req})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/resolve", c.baseURL)

	for {
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

		if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusServiceUnavailable {
			retryAfter := 30 * time.Second
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, parseErr := strconv.Atoi(ra); parseErr == nil {
					retryAfter = time.Duration(secs) * time.Second
				}
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			select {
			case <-time.After(retryAfter):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("idmsu service returned status %d: %s", resp.StatusCode, string(body))
	}

	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var results []ResolveResponseItem
	if err = json.Unmarshal(respBody, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return results, nil
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
