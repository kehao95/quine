package llm

import (
	"errors"
	"io"
	"net"
	"strings"
)

var transientInferenceDecodeSubstrings = []string{
	"parsing SSE stream: no response.completed event found in SSE stream",
	"parsing SSE stream: scanning SSE stream:",
}

var transientInferenceTransportSubstrings = []string{
	"connection reset by peer",
	"unexpected eof",
	"broken pipe",
	"use of closed network connection",
}

func isTransientInferenceDecodeError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, part := range transientInferenceDecodeSubstrings {
		if strings.Contains(msg, part) {
			return true
		}
	}
	return false
}

func isTransientInferenceTransportError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}
	msg := strings.ToLower(err.Error())
	for _, part := range transientInferenceTransportSubstrings {
		if strings.Contains(msg, part) {
			return true
		}
	}
	return false
}
