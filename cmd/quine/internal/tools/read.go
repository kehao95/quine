package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kehao95/quine/cmd/quine/internal/tape"
)

// countingReader wraps an io.Reader and counts bytes read.
// This is placed BEFORE bufio so we track raw reads from the source.
type countingReader struct {
	r     io.Reader
	count int64 // atomic for safety
}

func (c *countingReader) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	atomic.AddInt64(&c.count, int64(n))
	return
}

func (c *countingReader) BytesRead() int64 {
	return atomic.LoadInt64(&c.count)
}

// ReadExecutor reads lines from stdin on demand (streaming input).
type ReadExecutor struct {
	reader        *bufio.Reader
	counter       *countingReader // tracks raw bytes from source
	eof           bool
	Timeout       time.Duration // Default timeout for blocking reads
	MaxLines      int           // Maximum lines per read (0 = unlimited, default 500)
	initialOffset int64         // Starting offset (from previous exec)
}

// NewReadExecutor creates a ReadExecutor wrapping the given io.Reader.
// The reader is typically os.Stdin, kept open for streaming reads.
func NewReadExecutor(r io.Reader, defaultTimeout time.Duration) *ReadExecutor {
	counter := &countingReader{r: r}
	return &ReadExecutor{
		counter:  counter,
		reader:   bufio.NewReader(counter),
		Timeout:  defaultTimeout,
		MaxLines: 500, // Default cap to prevent context explosion
	}
}

// NewReadExecutorWithOffset creates a ReadExecutor with a known starting offset.
// Use this when resuming after exec to track total bytes consumed.
func NewReadExecutorWithOffset(r io.Reader, defaultTimeout time.Duration, initialOffset int64) *ReadExecutor {
	counter := &countingReader{r: r}
	return &ReadExecutor{
		counter:       counter,
		reader:        bufio.NewReader(counter),
		Timeout:       defaultTimeout,
		MaxLines:      500,
		initialOffset: initialOffset,
	}
}

// BytesConsumed returns the logical bytes consumed from the stream.
// This accounts for bufio's internal buffer, returning the position
// of the last byte we've actually returned to the caller (not the raw
// kernel position which may be ahead due to bufio's read-ahead).
//
// For seekable inputs (files), this is the position to seek to after exec.
// For non-seekable inputs (pipes), this is best-effort — some data in
// bufio's buffer may be lost if exec happens mid-buffer.
func (r *ReadExecutor) BytesConsumed() int64 {
	rawBytesRead := r.counter.BytesRead()
	bufferedBytes := int64(r.reader.Buffered())
	logicalPosition := rawBytesRead - bufferedBytes
	return r.initialOffset + logicalPosition
}

// ReadRequest represents parsed arguments for the read tool.
type ReadRequest struct {
	Lines   int `json:"lines"`   // Number of lines to read (default 1, 0 = read to EOF)
	Timeout int `json:"timeout"` // Timeout in seconds (0 = use default)
}

// ParseReadArgs parses the read tool arguments.
func ParseReadArgs(args map[string]any) (ReadRequest, error) {
	req := ReadRequest{Lines: 1} // Default: read 1 line

	data, err := json.Marshal(args)
	if err != nil {
		return req, fmt.Errorf("marshal args: %w", err)
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return req, fmt.Errorf("unmarshal args: %w", err)
	}

	if req.Lines < 0 {
		return req, fmt.Errorf("lines must be >= 0 (got %d)", req.Lines)
	}

	return req, nil
}

// Execute reads from stdin and returns a ToolResult.
//
// Behavior:
//   - lines=1 (default): Read one line, return it
//   - lines=N: Read N lines, return them joined with newlines
//   - lines=0: Read up to MaxLines (default 500), NOT all content
//
// The result includes metadata:
//   - [LINES READ] N - number of lines actually read
//   - [EOF] true/false - whether EOF was reached
//   - [TRUNCATED] true - if output was capped at MaxLines (more data available)
//   - [CONTENT] the actual data
func (r *ReadExecutor) Execute(toolID string, req ReadRequest) tape.ToolResult {
	if r.eof {
		return tape.ToolResult{
			ToolID:  toolID,
			Content: "[LINES READ] 0\n[EOF] true\n[CONTENT]\n",
			IsError: false,
		}
	}

	// Determine timeout
	timeout := r.Timeout
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var lines []string
	var eofReached bool
	var truncated bool

	// Determine how many lines to read
	linesToRead := req.Lines
	if linesToRead == 0 {
		// lines=0 means "read a batch", capped at MaxLines
		linesToRead = r.MaxLines
	}
	// Also cap explicit requests at MaxLines
	if r.MaxLines > 0 && linesToRead > r.MaxLines {
		linesToRead = r.MaxLines
		truncated = true
	}

	// Read N lines
	for i := 0; i < linesToRead; i++ {
		line, err := r.readLineWithContext(ctx)
		if err == io.EOF {
			r.eof = true
			eofReached = true
			break
		}
		if err != nil {
			return tape.ToolResult{
				ToolID:  toolID,
				Content: fmt.Sprintf("[ERROR] %v", err),
				IsError: true,
			}
		}
		lines = append(lines, line)
	}

	// If we read exactly linesToRead and it was capped, there might be more
	if len(lines) == linesToRead && !eofReached && r.MaxLines > 0 {
		truncated = true
	}

	content := strings.Join(lines, "\n")

	// Build result with metadata
	var result string
	if truncated {
		result = fmt.Sprintf("[LINES READ] %d\n[EOF] %v\n[TRUNCATED] true — more data available, call read again or use exec to reset context\n[CONTENT]\n%s", len(lines), eofReached, content)
	} else {
		result = fmt.Sprintf("[LINES READ] %d\n[EOF] %v\n[CONTENT]\n%s", len(lines), eofReached, content)
	}

	return tape.ToolResult{
		ToolID:  toolID,
		Content: result,
		IsError: false,
	}
}

// readLineWithContext reads a single line, respecting context cancellation.
func (r *ReadExecutor) readLineWithContext(ctx context.Context) (string, error) {
	type result struct {
		line string
		err  error
	}

	ch := make(chan result, 1)
	go func() {
		line, err := r.reader.ReadString('\n')
		// Trim the trailing newline
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		ch <- result{line, err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-ch:
		return res.line, res.err
	}
}

// EOF returns whether the underlying reader has reached EOF.
func (r *ReadExecutor) EOF() bool {
	return r.eof
}
