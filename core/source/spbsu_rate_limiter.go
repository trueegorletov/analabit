package source

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"log/slog"
)

type SpbsuRateLimiter struct {
	mu              sync.Mutex
	pausedUntil     time.Time
	backoffLevel    int
	backoffs        []time.Duration
	lastPauseReason string
}

var GlobalSpbsuRateLimiter = &SpbsuRateLimiter{
	backoffs: []time.Duration{
		10 * time.Second, 12 * time.Second, 15 * time.Second, 12 * time.Second, 15 * time.Second, 15 * time.Second, 30 * time.Second,
	},
}

func (rl *SpbsuRateLimiter) Wait(ctx context.Context) error {
	for {
		rl.mu.Lock()
		now := time.Now()
		if now.After(rl.pausedUntil) {
			rl.mu.Unlock()
			return nil
		}
		delay := rl.pausedUntil.Sub(now)
		rl.mu.Unlock()

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (rl *SpbsuRateLimiter) Handle429() {
	rl.mu.Lock()
	now := time.Now()
	backoff := rl.backoffs[rl.backoffLevel%len(rl.backoffs)]
	newPause := now.Add(backoff)
	if newPause.After(rl.pausedUntil) {
		rl.pausedUntil = newPause
	}
	rl.backoffLevel++
	rl.lastPauseReason = "429 error"
	rl.mu.Unlock()
	slog.Info("SPbSU paused for " + backoff.String() + " at level " + strconv.Itoa(rl.backoffLevel))
}

func (rl *SpbsuRateLimiter) Reset() {
	rl.mu.Lock()
	rl.backoffLevel = 0
	rl.pausedUntil = time.Time{}
	rl.mu.Unlock()
}

func (rl *SpbsuRateLimiter) MakeRequest(ctx context.Context, url string) (*http.Response, error) {
	var attempt int
	for {
		attempt++
		if err := rl.Wait(ctx); err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		client := &http.Client{Timeout: 0} // No timeout, rely on context
		resp, err := client.Do(req)
		if err != nil {
			slog.Warn("SPbSU request error, retrying", "attempt", attempt, "url", url, "error", err)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		if resp.StatusCode == 429 {
			resp.Body.Close()
			rl.Handle429()
			slog.Warn("SPbSU 429 error, retrying", "attempt", attempt, "url", url)
			continue
		}
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			resp.Body.Close()
			slog.Warn("SPbSU server error, retrying", "attempt", attempt, "url", url, "status", resp.StatusCode)
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		rl.Reset()
		return resp, nil
	}
}
