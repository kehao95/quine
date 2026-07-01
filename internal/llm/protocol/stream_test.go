package protocol

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kehao95/quine/internal/tape"
)

// T-SSE-INC: DecodeStream must surface deltas as they arrive on the wire, not
// only after the whole stream has been buffered. We prove this by holding the
// stream open (writer blocked) after the first delta and asserting onDelta has
// already fired.
func TestOpenAIResponsesDecodeStream_EmitsDeltasIncrementally(t *testing.T) {
	pr, pw := io.Pipe()

	firstChunk := "event: response.output_text.delta\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"Hello, \"}\n" +
		"\n"
	restChunks := "event: response.output_text.delta\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"world\"}\n" +
		"\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"output\":[{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hello, world\"}]}],\"usage\":{\"input_tokens\":3,\"output_tokens\":2}}}\n" +
		"\n"

	var (
		mu        sync.Mutex
		deltas    []string
		firstSeen = make(chan struct{})
		firstOnce sync.Once
	)
	onDelta := func(d StreamDelta) {
		if d.Kind != StreamDeltaText {
			return
		}
		mu.Lock()
		deltas = append(deltas, d.Text)
		mu.Unlock()
		firstOnce.Do(func() { close(firstSeen) })
	}

	type result struct {
		msg   tape.Message
		usage Usage
		err   error
	}
	done := make(chan result, 1)
	go func() {
		msg, usage, err := (&OpenAIResponsesProtocol{}).DecodeStream(pr, onDelta)
		done <- result{msg, usage, err}
	}()

	unblock := make(chan struct{})
	go func() {
		// Write the first event and the blank line that flushes it, then hold
		// the stream open until the test releases us.
		if _, err := io.WriteString(pw, firstChunk); err != nil {
			return
		}
		<-unblock
		_, _ = io.WriteString(pw, restChunks)
		_ = pw.Close()
	}()

	// The first delta must be observable while the writer is still blocked, i.e.
	// before the terminal response.completed event has been written.
	select {
	case <-firstSeen:
	case <-time.After(2 * time.Second):
		close(unblock)
		<-done
		t.Fatal("first delta not emitted before the stream completed (buffered, not incremental)")
	}

	close(unblock)

	var res result
	select {
	case res = <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("DecodeStream did not return after stream close")
	}
	if res.err != nil {
		t.Fatalf("DecodeStream() error = %v", res.err)
	}

	mu.Lock()
	got := strings.Join(deltas, "")
	mu.Unlock()
	if got != "Hello, world" {
		t.Fatalf("concatenated deltas = %q, want %q", got, "Hello, world")
	}
	if res.msg.Content != "Hello, world" {
		t.Fatalf("final content = %q, want %q", res.msg.Content, "Hello, world")
	}
	if res.usage.InputTokens != 3 || res.usage.OutputTokens != 2 {
		t.Fatalf("final usage = %+v", res.usage)
	}
}

// T-REPRO: feeding the same recorded SSE body through the buffered DecodeResponse
// and the streaming DecodeStream must crystallize an identical final cell. The
// only legitimate difference is the wall-clock Timestamp, which we zero out.
func TestOpenAIResponsesDecodeStream_FinalMatchesBuffered(t *testing.T) {
	cases := map[string]string{
		"completed_with_inline_output": strings.Join([]string{
			"event: response.reasoning_summary_text.delta",
			"data: {\"type\":\"response.reasoning_summary_text.delta\",\"item_id\":\"rs_1\",\"summary_index\":0,\"delta\":\"think \"}",
			"",
			"event: response.reasoning_summary_text.delta",
			"data: {\"type\":\"response.reasoning_summary_text.delta\",\"item_id\":\"rs_1\",\"summary_index\":0,\"delta\":\"hard\"}",
			"",
			"event: response.output_text.delta",
			"data: {\"type\":\"response.output_text.delta\",\"output_index\":0,\"content_index\":0,\"delta\":\"Hel\"}",
			"",
			"event: response.output_text.delta",
			"data: {\"type\":\"response.output_text.delta\",\"output_index\":0,\"content_index\":0,\"delta\":\"lo\"}",
			"",
			"event: response.completed",
			"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"output\":[{\"id\":\"rs_1\",\"type\":\"reasoning\"},{\"type\":\"message\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hello\"}]}],\"usage\":{\"input_tokens\":11,\"output_tokens\":4}}}",
			"",
		}, "\n"),
		// Empty completed output reconstructed from output_item.done plus merged
		// reasoning summary deltas — exercises the synthesize/merge paths.
		"reconstructed_from_output_item_done": strings.Join([]string{
			"event: response.reasoning_summary_text.delta",
			"data: {\"type\":\"response.reasoning_summary_text.delta\",\"item_id\":\"rs_9\",\"summary_index\":0,\"delta\":\"summary text\"}",
			"",
			"event: response.output_item.done",
			"data: {\"type\":\"response.output_item.done\",\"output_index\":1,\"item\":{\"id\":\"msg_9\",\"type\":\"message\",\"status\":\"completed\",\"role\":\"assistant\",\"content\":[{\"type\":\"output_text\",\"text\":\"Reconstructed\"}]}}",
			"",
			"event: response.completed",
			"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_9\",\"output\":[{\"id\":\"rs_9\",\"type\":\"reasoning\"}],\"usage\":{\"input_tokens\":7,\"output_tokens\":13}}}",
			"",
		}, "\n"),
	}

	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			raw := []byte(body)

			bufMsg, bufUsage, err := (&OpenAIResponsesProtocol{}).DecodeResponse(raw)
			if err != nil {
				t.Fatalf("DecodeResponse() error = %v", err)
			}
			strMsg, strUsage, err := (&OpenAIResponsesProtocol{}).DecodeStream(bytes.NewReader(raw), nil)
			if err != nil {
				t.Fatalf("DecodeStream() error = %v", err)
			}

			// Timestamp is wall-clock and intentionally differs between calls.
			bufMsg.Timestamp = 0
			strMsg.Timestamp = 0

			if !reflect.DeepEqual(bufMsg, strMsg) {
				t.Fatalf("final message mismatch:\n buffered = %#v\nstreaming = %#v", bufMsg, strMsg)
			}
			if bufUsage != strUsage {
				t.Fatalf("usage mismatch: buffered = %+v streaming = %+v", bufUsage, strUsage)
			}
		})
	}
}

// T-ERROR: a terminal failure event with no response.completed must propagate
// as a non-nil error carrying the vendor code/message.
func TestOpenAIResponsesDecodeStream_PropagatesTerminalFailure(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []string
	}{
		{
			name: "response.failed",
			body: strings.Join([]string{
				"event: response.output_text.delta",
				"data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"}",
				"",
				"event: response.failed",
				"data: {\"type\":\"response.failed\",\"response\":{\"error\":{\"code\":\"server_error\",\"message\":\"upstream exploded\"}}}",
				"",
			}, "\n"),
			want: []string{"response.failed", "server_error", "upstream exploded"},
		},
		{
			name: "error_event",
			body: strings.Join([]string{
				"event: error",
				"data: {\"type\":\"error\",\"error\":{\"type\":\"invalid_request\",\"code\":\"bad_input\",\"message\":\"malformed stream\"}}",
				"",
			}, "\n"),
			want: []string{"error", "bad_input", "malformed stream"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := (&OpenAIResponsesProtocol{}).DecodeStream(strings.NewReader(tc.body), func(StreamDelta) {})
			if err == nil {
				t.Fatal("DecodeStream() error = nil, want terminal failure")
			}
			for _, want := range tc.want {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error %q does not contain %q", err.Error(), want)
				}
			}
		})
	}
}
