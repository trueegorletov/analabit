package source

import (
	"context"
	"log/slog"
)

// MSUReceiverFactory creates MSU-specific receivers to avoid import cycles
type MSUReceiverFactory interface {
	CreateMSUReceiver(downstream DataReceiver) MSUReceiver
}

// MSUReceiver extends DataReceiver with finalization for ID resolution
type MSUReceiver interface {
	DataReceiver
	Finalize(ctx context.Context) error
}

// DefaultMSUReceiverFactory is a fallback that just returns the downstream receiver
type DefaultMSUReceiverFactory struct{}

func (f *DefaultMSUReceiverFactory) CreateMSUReceiver(downstream DataReceiver) MSUReceiver {
	slog.Warn("Using default MSU receiver factory - MSU ID resolution not available")
	return &defaultMSUReceiver{downstream: downstream}
}

// defaultMSUReceiver is a passthrough receiver for when MSU buffering is not available
type defaultMSUReceiver struct {
	downstream DataReceiver
}

func (r *defaultMSUReceiver) PutHeadingData(hd *HeadingData) {
	r.downstream.PutHeadingData(hd)
}

func (r *defaultMSUReceiver) PutApplicationData(ad *ApplicationData) {
	r.downstream.PutApplicationData(ad)
}

func (r *defaultMSUReceiver) Finalize(ctx context.Context) error {
	// No-op for default receiver
	return nil
}

// Global factory instance - can be overridden by MSU package
var MSUFactory MSUReceiverFactory = &DefaultMSUReceiverFactory{}

// SetMSUReceiverFactory allows the MSU package to register its factory
func SetMSUReceiverFactory(factory MSUReceiverFactory) {
	MSUFactory = factory
}