package protocol

import (
	"io"

	"github.com/kehao95/quine/internal/tape"
)

// StreamDeltaKind classifies an incremental generation delta emitted while a
// response is still being produced. Deltas are transient display signal only;
// the authoritative final cell is always the value returned by DecodeStream,
// derived from the same terminal payload the buffered DecodeResponse uses.
type StreamDeltaKind int

const (
	// StreamDeltaText is an incremental piece of assistant-visible output text.
	StreamDeltaText StreamDeltaKind = iota
	// StreamDeltaReasoning is an incremental piece of reasoning-summary text.
	StreamDeltaReasoning
	// StreamDeltaToolCall is an incremental piece of a tool call's arguments.
	StreamDeltaToolCall
)

// StreamDelta is one incremental generation event surfaced to a streaming
// consumer. It never feeds provider input or recovery — it exists purely so a
// live preview can render generation as it arrives.
type StreamDelta struct {
	Kind     StreamDeltaKind
	Text     string
	ToolID   string
	ToolName string
}

// StreamingProtocol is implemented by protocols that can decode an HTTP
// response body incrementally, invoking onDelta as each delta arrives and
// returning the authoritative final (message, usage) once the stream completes.
//
// The final cell MUST be derived from the same terminal payload the buffered
// DecodeResponse uses (shared final-assembly), so the streaming and buffered
// paths produce an equivalent crystallized cell (reproducibility boundary).
type StreamingProtocol interface {
	DecodeStream(r io.Reader, onDelta func(StreamDelta)) (tape.Message, Usage, error)
}
