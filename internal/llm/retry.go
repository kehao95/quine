package llm

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// ---------------------------------------------------------------------------
// Retry with exponential backoff + jitter
// ---------------------------------------------------------------------------

// retryWithBackoff executes fn with retry logic. The retry behaviour depends
// on the HTTP status code returned:
//
//   - 429  → up to maxRetries retries with exponential backoff + jitter
//   - 5xx  → up to 3 retries
//   - 401/403 → no retry, return immediately
//   - Network error → up to 3 retries
//   - Malformed / unexpected → retry once
func retryWithBackoff(maxRetries int, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := fn()

		if err != nil {
			// Network-level error → retry up to 3 times
			lastErr = err
			retryLimit := 3
			if attempt >= retryLimit {
				return nil, lastErr
			}
			logRetry(attempt+1, retryLimit, err.Error())
			backoffSleep(attempt)
			continue
		}

		switch {
		case resp.StatusCode == http.StatusOK:
			return resp, nil

		case resp.StatusCode == 429:
			// Rate limited → retry up to maxRetries times
			drainAndClose(resp)
			lastErr = fmt.Errorf("rate limited (429)")
			if attempt >= maxRetries {
				return fn()
			}
			logRetry(attempt+1, maxRetries, "rate limited (429)")
			backoffSleep(attempt)
			continue

		case resp.StatusCode == 401 || resp.StatusCode == 403:
			// Auth error → no retry
			return resp, nil

		case resp.StatusCode >= 500:
			// Server error → retry up to 3 times
			drainAndClose(resp)
			lastErr = fmt.Errorf("server error (%d)", resp.StatusCode)
			retryLimit := 3
			if attempt >= retryLimit {
				return fn()
			}
			logRetry(attempt+1, retryLimit, fmt.Sprintf("server error (%d)", resp.StatusCode))
			backoffSleep(attempt)
			continue

		default:
			// Unexpected status → retry once
			if attempt >= 1 {
				return resp, nil
			}
			drainAndClose(resp)
			lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
			logRetry(attempt+1, 1, fmt.Sprintf("unexpected status (%d)", resp.StatusCode))
			backoffSleep(attempt)
			continue
		}
	}

	if lastResp != nil {
		return lastResp, nil
	}
	return nil, lastErr
}

func logRetry(attempt, max int, reason string) {
	fmt.Fprintf(stderrWriter(), "quine: LLM retry %d/%d (%s)\n", attempt, max, reason)
}

func stderrWriter() io.Writer {
	return stderrOut
}

var stderrOut io.Writer = os.Stderr

// SetLogOutput redirects operational log messages to w instead of os.Stderr.
func SetLogOutput(w io.Writer) {
	if w == nil {
		stderrOut = io.Discard
		return
	}
	stderrOut = w
}

func backoffSleep(attempt int) {
	base := time.Duration(1<<uint(attempt)) * 500 * time.Millisecond
	jitter := time.Duration(rand.Int63n(int64(base / 2)))
	time.Sleep(base + jitter)
}

func drainAndClose(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
