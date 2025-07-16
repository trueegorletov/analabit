package flaresolverr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

// SessionInfo represents information about a FlareSolverr session
type SessionInfo struct {
	ID           string
	Domain       string
	CreatedAt    time.Time
	LastUsedAt   time.Time
	RequestCount int
	ActiveReqs   int
	Healthy      bool
	mutex        sync.RWMutex
}

// SessionPool manages sessions for a specific domain
type SessionPool struct {
	domain   string
	sessions []*SessionInfo
	mutex    sync.RWMutex
	maxSize  int
	client   *Client
}

// SessionManager manages session pools for different domains
type SessionManager struct {
	pools                 map[string]*SessionPool
	mutex                 sync.RWMutex
	client                *Client
	poolSize              int
	idleTimeout           time.Duration
	maxRequestsPerSession int
	healthCheckInterval   time.Duration
	shutdownChan          chan struct{}
	shutdownOnce          sync.Once
}

var (
	clientInstance *Client
	clientOnce     sync.Once
	clientError    error
	sessionManager *SessionManager
	sessionOnce    sync.Once
)

// GetClient returns a singleton FlareSolverr client instance
func GetClient() (*Client, error) {
	clientOnce.Do(func() {
		clientInstance, clientError = initializeClient()
	})
	return clientInstance, clientError
}

// GetSessionManager returns a singleton SessionManager instance
func GetSessionManager() (*SessionManager, error) {
	sessionOnce.Do(func() {
		client, err := GetClient()
		if err != nil {
			clientError = err
			return
		}
		sessionManager = initializeSessionManager(client)
	})
	return sessionManager, clientError
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

// initializeSessionManager creates and configures a new SessionManager
func initializeSessionManager(client *Client) *SessionManager {
	// Load configuration from environment variables
	poolSize := 8 // Default pool size
	if envPoolSize := os.Getenv("FLARESOLVERR_SESSION_POOL_SIZE"); envPoolSize != "" {
		if parsed, err := strconv.Atoi(envPoolSize); err == nil && parsed > 0 {
			poolSize = parsed
		}
	}

	idleTimeout := 5 * time.Minute // Default idle timeout
	if envTimeout := os.Getenv("FLARESOLVERR_SESSION_IDLE_TIMEOUT_MINUTES"); envTimeout != "" {
		if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
			idleTimeout = time.Duration(parsed) * time.Minute
		}
	}

	maxRequestsPerSession := 80 // Default max requests per session
	if envMaxReqs := os.Getenv("FLARESOLVERR_SESSION_MAX_REQUESTS_PER_SESSION"); envMaxReqs != "" {
		if parsed, err := strconv.Atoi(envMaxReqs); err == nil && parsed > 0 {
			maxRequestsPerSession = parsed
		}
	}

	healthCheckInterval := 30 * time.Second // Default health check interval
	if envHealthCheck := os.Getenv("FLARESOLVERR_SESSION_HEALTH_CHECK_INTERVAL_SECONDS"); envHealthCheck != "" {
		if parsed, err := strconv.Atoi(envHealthCheck); err == nil && parsed > 0 {
			healthCheckInterval = time.Duration(parsed) * time.Second
		}
	}

	sm := &SessionManager{
		pools:                 make(map[string]*SessionPool),
		client:                client,
		poolSize:              poolSize,
		idleTimeout:           idleTimeout,
		maxRequestsPerSession: maxRequestsPerSession,
		healthCheckInterval:   healthCheckInterval,
		shutdownChan:          make(chan struct{}),
	}

	// Start background cleanup and health check routine
	go sm.backgroundMaintenance()

	return sm
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

// GetWithDomain performs a GET request using domain-specific session management
// This function automatically manages sessions for the domain extracted from the URL
func GetWithDomain(url string, headers map[string]string) (*GetResponse, error) {
	// Extract domain from URL
	domain, err := extractDomain(url)
	if err != nil {
		return &GetResponse{Error: err}, fmt.Errorf("failed to extract domain from URL: %w", err)
	}

	// Get session manager
	sessionManager, err := GetSessionManager()
	if err != nil {
		// Fallback to sessionless request
		return SafeGetWithHeaders(url, headers)
	}

	// Get session for domain
	session, err := sessionManager.GetSessionForDomain(domain)
	if err != nil {
		// Fallback to sessionless request
		return SafeGetWithHeaders(url, headers)
	}

	// Make request with session
	client, err := GetClient()
	if err != nil {
		return &GetResponse{Error: err}, fmt.Errorf("FlareSolverr unavailable: %w", err)
	}

	resp, err := client.GetWithSession(url, session.ID, headers)

	// Update session stats
	session.mutex.Lock()
	session.RequestCount++
	session.mutex.Unlock()

	// Release session back to pool
	sessionManager.ReleaseSession(session)

	if err != nil {
		// If session request failed, try to recreate session or fallback
		if strings.Contains(err.Error(), "session") {
			// Session might be invalid, mark as unhealthy
			session.mutex.Lock()
			session.Healthy = false
			session.mutex.Unlock()
		}
		return resp, err
	}

	return resp, nil
}

// SafeGetWithDomain performs a GET request with domain-specific session management and graceful fallback
func SafeGetWithDomain(url string, headers map[string]string) (*GetResponse, error) {
	resp, err := GetWithDomain(url, headers)
	if err != nil && IsFlareSolverrError(err) {
		return &GetResponse{Error: err}, fmt.Errorf("FlareSolverr unavailable: %w", err)
	}
	return resp, err
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

// SessionManager methods

// extractDomain extracts the domain from a URL
func extractDomain(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	return strings.ToLower(parsedURL.Host), nil
}

// getOrCreatePool gets or creates a session pool for a domain
func (sm *SessionManager) getOrCreatePool(domain string) *SessionPool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if pool, exists := sm.pools[domain]; exists {
		return pool
	}

	pool := &SessionPool{
		domain:   domain,
		sessions: make([]*SessionInfo, 0, sm.poolSize),
		maxSize:  sm.poolSize,
		client:   sm.client,
	}
	sm.pools[domain] = pool
	return pool
}

// GetSessionForDomain gets an available session for a domain
func (sm *SessionManager) GetSessionForDomain(domain string) (*SessionInfo, error) {
	pool := sm.getOrCreatePool(domain)
	return pool.getAvailableSession()
}

// ReleaseSession releases a session back to the pool
func (sm *SessionManager) ReleaseSession(session *SessionInfo) {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	session.ActiveReqs--
	session.LastUsedAt = time.Now()
}

// Shutdown gracefully shuts down the session manager
func (sm *SessionManager) Shutdown(ctx context.Context) error {
	sm.shutdownOnce.Do(func() {
		close(sm.shutdownChan)
	})

	// Wait for background maintenance to stop
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		// Continue with cleanup
	}

	// Destroy all sessions
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, pool := range sm.pools {
		pool.destroyAllSessions()
	}

	return nil
}

// backgroundMaintenance runs background cleanup and health checks
func (sm *SessionManager) backgroundMaintenance() {
	ticker := time.NewTicker(sm.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.shutdownChan:
			return
		case <-ticker.C:
			sm.cleanupIdleSessions()
			sm.healthCheckSessions()
		}
	}
}

// cleanupIdleSessions removes idle sessions that exceed the timeout
func (sm *SessionManager) cleanupIdleSessions() {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, pool := range sm.pools {
		pool.cleanupIdleSessions(sm.idleTimeout)
	}
}

// healthCheckSessions performs health checks on sessions
func (sm *SessionManager) healthCheckSessions() {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for _, pool := range sm.pools {
		pool.healthCheckSessions()
	}
}

// SessionPool methods

// getAvailableSession gets an available session from the pool
func (sp *SessionPool) getAvailableSession() (*SessionInfo, error) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	// Find session with least active requests
	var bestSession *SessionInfo
	minActiveReqs := int(^uint(0) >> 1) // Max int

	for _, session := range sp.sessions {
		session.mutex.RLock()
		if session.Healthy && session.ActiveReqs < minActiveReqs {
			bestSession = session
			minActiveReqs = session.ActiveReqs
		}
		session.mutex.RUnlock()
	}

	// If we found a good session, use it
	if bestSession != nil {
		bestSession.mutex.Lock()
		bestSession.ActiveReqs++
		bestSession.mutex.Unlock()
		return bestSession, nil
	}

	// No available session, create a new one if pool not full
	if len(sp.sessions) < sp.maxSize {
		session, err := sp.createNewSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create new session: %w", err)
		}
		session.ActiveReqs = 1
		sp.sessions = append(sp.sessions, session)
		return session, nil
	}

	// Pool is full, wait for least busy session
	if bestSession == nil && len(sp.sessions) > 0 {
		bestSession = sp.sessions[0]
		for _, session := range sp.sessions[1:] {
			session.mutex.RLock()
			bestSession.mutex.RLock()
			if session.ActiveReqs < bestSession.ActiveReqs {
				bestSession.mutex.RUnlock()
				bestSession = session
				session.mutex.RUnlock()
			} else {
				bestSession.mutex.RUnlock()
				session.mutex.RUnlock()
			}
		}
		bestSession.mutex.Lock()
		bestSession.ActiveReqs++
		bestSession.mutex.Unlock()
		return bestSession, nil
	}

	return nil, fmt.Errorf("no sessions available for domain %s", sp.domain)
}

// createNewSession creates a new session for the pool
func (sp *SessionPool) createNewSession() (*SessionInfo, error) {
	sessionID, err := sp.client.CreateSession()
	if err != nil {
		return nil, err
	}

	return &SessionInfo{
		ID:           sessionID,
		Domain:       sp.domain,
		CreatedAt:    time.Now(),
		LastUsedAt:   time.Now(),
		RequestCount: 0,
		ActiveReqs:   0,
		Healthy:      true,
	}, nil
}

// cleanupIdleSessions removes sessions that have been idle too long
func (sp *SessionPool) cleanupIdleSessions(idleTimeout time.Duration) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	now := time.Now()
	var activeSessions []*SessionInfo

	for _, session := range sp.sessions {
		session.mutex.RLock()
		isIdle := session.ActiveReqs == 0 && now.Sub(session.LastUsedAt) > idleTimeout
		session.mutex.RUnlock()

		if isIdle {
			// Destroy idle session
			sp.client.DestroySession(session.ID)
		} else {
			activeSessions = append(activeSessions, session)
		}
	}

	sp.sessions = activeSessions
}

// healthCheckSessions performs health checks on pool sessions
func (sp *SessionPool) healthCheckSessions() {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	for _, session := range sp.sessions {
		go sp.checkSessionHealth(session)
	}
}

// checkSessionHealth checks if a session is still healthy
func (sp *SessionPool) checkSessionHealth(session *SessionInfo) {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	// Simple health check - if session has too many requests, mark as unhealthy
	if session.RequestCount > 100 { // Configurable threshold
		session.Healthy = false
	}
}

// destroyAllSessions destroys all sessions in the pool
func (sp *SessionPool) destroyAllSessions() {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	for _, session := range sp.sessions {
		sp.client.DestroySession(session.ID)
	}
	sp.sessions = nil
}
