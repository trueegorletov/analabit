package utils

import (
	"analabit/core/ent"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
)

// PrepareStudentID Removes all non-numeric characters from the student ID and returns it
// being exactly 13 characters long. If less, it pads with zeros at the start
func PrepareStudentID(rawId string) (string, error) {
	var id string
	for _, c := range rawId {
		if c >= '0' && c <= '9' {
			id += string(c)
		}
	}

	if len(id) > 13 {
		return "", fmt.Errorf("student ID %q is longer than 13 characters after sanitization: %q", rawId, id)
	}

	for len(id) < 13 {
		id = "0" + id
	}
	return id, nil
}

// GenerateHeadingCode creates a SHA256 hash of the name and returns its hex representation.
func GenerateHeadingCode(name string) string {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	return hex.EncodeToString(hasher.Sum(nil))
}

// MustParseURL parses a raw URL string and returns a pointer to a url.URL.
// It logs.Fatalf if parsing fails. If rawURL is empty, it returns nil.
func MustParseURL(rawURL string) *url.URL {
	if rawURL == "" {
		return nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("failed to parse URL %q: %v", rawURL, err)
	}
	return u
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
