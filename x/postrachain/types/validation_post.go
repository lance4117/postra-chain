package types

import (
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	maxTitleLen      = 140
	maxContentURILen = 1024
)

var contentHashPattern = regexp.MustCompile(`^(?i)(sha256:)?[0-9a-f]{64}$`)

func ValidatePostFields(title, contentURI, contentHash string) error {
	if err := ValidateTitle(title); err != nil {
		return err
	}
	if err := ValidateContentURI(contentURI); err != nil {
		return err
	}
	if err := ValidateContentHash(contentHash); err != nil {
		return err
	}
	return nil
}

func ValidateTitle(title string) error {
	titleLen := utf8.RuneCountInString(title)
	if titleLen == 0 || titleLen > maxTitleLen {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(title) == "" {
		return ErrInvalidTitle
	}
	return nil
}

func ValidateContentURI(contentURI string) error {
	if contentURI == "" || len(contentURI) > maxContentURILen {
		return ErrInvalidContentURI
	}
	parsed, err := url.Parse(contentURI)
	if err != nil || parsed.Scheme == "" {
		return ErrInvalidContentURI
	}
	switch strings.ToLower(parsed.Scheme) {
	case "http", "https", "ipfs":
		return nil
	default:
		return ErrInvalidContentURI
	}
}

func ValidateContentHash(contentHash string) error {
	if contentHash == "" || !contentHashPattern.MatchString(contentHash) {
		return ErrInvalidContentHash
	}
	return nil
}
