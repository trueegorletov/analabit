package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trueegorletov/analabit/core/idresolver"
)

// Handler aggregates all HTTP handlers of idmsu.
// It is designed to be thin; business logic resides in the resolver package.

type Handler struct {
	res        idresolver.StudentIDResolver
	readyCheck func() bool
}

func NewHandler(res idresolver.StudentIDResolver, readyCheck func() bool) *Handler {
	return &Handler{res: res, readyCheck: readyCheck}
}

// ResolveBatch handles POST /api/v1/resolve.
// It supports the same contract as the legacy service but is currently a stub.
func (h *Handler) ResolveBatch(c *gin.Context) {
	if h.readyCheck != nil && !h.readyCheck() {
		c.Header("Retry-After", "30")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service not ready"})
		return
	}

	// Read entire body once
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body", "details": err.Error()})
		return
	}

	// Attempt legacy format {"items": [...]}
	var legacyWrapper struct {
		Items []idresolver.ResolveRequestItem `json:"items"`
	}
	if err := json.Unmarshal(raw, &legacyWrapper); err == nil && len(legacyWrapper.Items) > 0 {
		h.processResolve(c, legacyWrapper.Items)
		return
	}

	// Attempt direct array format
	var req []idresolver.ResolveRequestItem
	if err := json.Unmarshal(raw, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	h.processResolve(c, req)
}

// processResolve handles the actual ID resolution using the sophisticated algorithm
func (h *Handler) processResolve(c *gin.Context, req []idresolver.ResolveRequestItem) {
	// Call the resolver to perform the actual matching
	resp, err := h.res.ResolveBatch(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "resolution failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Health returns liveness status.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready indicates readiness; currently always false until background cache populated.
func (h *Handler) Ready(c *gin.Context) {
	if h.readyCheck != nil && h.readyCheck() {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
		return
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{"status": "initializing"})
}

// Wait provides asynchronous readiness wait endpoint following 202 + Retry-After semantics.
func (h *Handler) Wait(c *gin.Context) {
	if h.readyCheck != nil && h.readyCheck() {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
		return
	}
	c.Header("Retry-After", "30")
	c.JSON(http.StatusAccepted, gin.H{"status": "processing"})
}
