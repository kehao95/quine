package protocol

import "errors"

// Sentinel errors for callers to match with errors.Is.
var (
	ErrAuth            = errors.New("authentication failed")
	ErrContextOverflow = errors.New("context window exceeded")
)
