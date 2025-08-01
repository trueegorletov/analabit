package idresolver

import (
	"context"
	"fmt"
	"log/slog"
)

// StubResolver is a temporary implementation of StudentIDResolver for testing
// This will be replaced by the actual idmsu service client
type StubResolver struct{}

// NewStubResolver creates a new stub resolver
func NewStubResolver() *StubResolver {
	return &StubResolver{}
}

// ResolveBatch implements StudentIDResolver interface with fallback logic
func (s *StubResolver) ResolveBatch(ctx context.Context, req []ResolveRequestItem) ([]ResolveResponseItem, error) {
	slog.Warn("Using stub resolver - idmsu service not implemented yet", "requestCount", len(req))
	
	response := make([]ResolveResponseItem, len(req))
	for i, item := range req {
		// Create fallback canonical ID: MSU- prefix + padded internal ID
		fallbackID := fmt.Sprintf("MSU-%s", item.InternalID)
		// Pad to 13 characters
		for len(fallbackID) < 13 {
			fallbackID = "0" + fallbackID
		}
		// If too long, truncate from beginning but preserve some internal ID
		if len(fallbackID) > 13 {
			fallbackID = fallbackID[len(fallbackID)-13:]
		}
		
		response[i] = ResolveResponseItem{
			InternalID:  item.InternalID,
			CanonicalID: fallbackID,
			Confidence:  0.0, // Low confidence since this is a stub
		}
		
		slog.Debug("Stub resolver created fallback ID", 
			"internalID", item.InternalID, 
			"canonicalID", fallbackID,
			"appCount", len(item.Apps))
	}
	
	return response, nil
}