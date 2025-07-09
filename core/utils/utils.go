package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strings"
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

// PrettifyStudentID removes leading zeros from a student ID stored in the database
// Example: "0000004272036" -> "4272036", "01012345" -> "1012345"
// Always preserves at least one digit, so "0000000000000" -> "0"
func PrettifyStudentID(storedId string) string {
	if storedId == "" {
		return storedId
	}

	// Use TrimLeft to remove leading zeros - this is optimized and zero-copy for prefix removal
	prettified := strings.TrimLeft(storedId, "0")

	// Ensure at least one digit remains
	if prettified == "" {
		return "0"
	}

	return prettified
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
