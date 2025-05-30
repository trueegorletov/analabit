package source

// HeadingSource is an interface for loading heading and application data.
// Implementations should send data to the provided channels.
// It is the responsibility of the caller to close the channels after LoadTo returns or if an error occurs.
type HeadingSource interface {
	LoadTo(headings chan<- HeadingData, applications chan<- ApplicationData) error
}
