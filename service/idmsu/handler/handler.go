package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trueegorletov/analabit/core/idresolver"
	"github.com/trueegorletov/analabit/service/idmsu/resolver"
)

type Handler struct {
	resolver resolver.MSUResolver
}

type ResolveBatchRequest struct {
	Items []idresolver.ResolveRequestItem `json:"items"`
}

type ResolveBatchResponse struct {
	Results []idresolver.ResolveResponseItem `json:"results"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type ReadinessResponse struct {
	Ready     bool      `json:"ready"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func NewHandler(resolver resolver.MSUResolver) *Handler {
	return &Handler{
		resolver: resolver,
	}
}

func (h *Handler) ResolveBatch(c *gin.Context) {
    if !h.resolver.HasRecentData(c.Request.Context()) {
        c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service is initializing - please wait"})
        return
    }
	var req ResolveBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items to resolve"})
		return
	}

	if len(req.Items) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many items, maximum 1000 allowed"})
		return
	}

	results, err := h.resolver.ResolveBatch(c.Request.Context(), req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve IDs", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Ready(c *gin.Context) {
	hasRecentData := h.resolver.HasRecentData(c.Request.Context())
	
	var response ReadinessResponse
	var statusCode int
	
	if hasRecentData {
		response = ReadinessResponse{
			Ready:     true,
			Timestamp: time.Now(),
			Message:   "Service is ready with recent Gosuslugi data",
		}
		statusCode = http.StatusOK
	} else {
		response = ReadinessResponse{
			Ready:     false,
			Timestamp: time.Now(),
			Message:   "Service is not ready - no recent Gosuslugi data available",
		}
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}