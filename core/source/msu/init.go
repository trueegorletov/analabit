package msu

import (
	"github.com/trueegorletov/analabit/core/idresolver"
)

func init() {
	// Always use the real idmsu service client
	// The MSU buffered receiver has fallback logic for when ID resolution fails
	resolver := idresolver.NewIDMSUClient()
	
	// Register the MSU receiver factory with the resolver
	RegisterFactory(resolver)
}