package tools

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/creack/pty"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	defaultInteractiveCols = 80
	defaultInteractiveRows = 24
)

type interactiveMeta struct {
	Rows       int    `json:"rows"`
	Cols       int    `json:"cols"`
	CursorRow  int    `json:"cursor_row"`
	CursorCol  int    `json:"cursor_col"`
	AltScreen  bool   `json:"alt_screen"`
	Generation uint64 `json:"generation"`
	UpdatedAt  string `json:"updated_at"`
}

type interactiveCell struct {
	Ch      rune
	FG      int
	BG      int
	Bold    bool
	Inverse bool
}

type interactiveBuffer struct {
	rows int
	cols int
	grid [][]interactiveCell
}

type interactiveScreen struct {
	rows int
	cols int

	cursorRow int
	cursorCol int
	savedRow  int
	savedCol  int

	mainBuf *interactiveBuffer
	altBuf  *interactiveBuffer

	altScreen bool

	fg      int
	bg      int
	bold    bool
	inverse bool

	escState   string
	csiParams  strings.Builder
	csiPrivate bool
}

type interactiveState struct {
	ptyFile     *os.File
	eventsFile  *os.File
	inputPath   string
	screenPath  string
	metaPath    string
	pngPath     string
	winsizePath string
	eventsPath  string
	eventsHex   *os.File
	inputLog    *os.File

	mu          sync.Mutex
	screen      *interactiveScreen
	generation  uint64
	lastWinsize string

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func newInteractiveBuffer(rows, cols int) *interactiveBuffer {
	buf := &interactiveBuffer{
		rows: rows,
		cols: cols,
		grid: make([][]interactiveCell, rows),
	}
	for r := 0; r < rows; r++ {
		buf.grid[r] = make([]interactiveCell, cols)
		for c := 0; c < cols; c++ {
			buf.grid[r][c] = blankCell()
		}
	}
	return buf
}

func blankCell() interactiveCell {
	return interactiveCell{Ch: ' ', FG: 7, BG: 0}
}

func newInteractiveScreen(rows, cols int) *interactiveScreen {
	return &interactiveScreen{
		rows:    rows,
		cols:    cols,
		mainBuf: newInteractiveBuffer(rows, cols),
		altBuf:  newInteractiveBuffer(rows, cols),
		fg:      7,
		bg:      0,
	}
}

func (s *interactiveScreen) currentBuffer() *interactiveBuffer {
	if s.altScreen {
		return s.altBuf
	}
	return s.mainBuf
}

func (s *interactiveScreen) resetAttrs() {
	s.fg = 7
	s.bg = 0
	s.bold = false
	s.inverse = false
}

func (s *interactiveScreen) reset() {
	s.cursorRow = 0
	s.cursorCol = 0
	s.savedRow = 0
	s.savedCol = 0
	s.altScreen = false
	s.mainBuf = newInteractiveBuffer(s.rows, s.cols)
	s.altBuf = newInteractiveBuffer(s.rows, s.cols)
	s.resetAttrs()
	s.escState = ""
	s.csiParams.Reset()
	s.csiPrivate = false
}

func (s *interactiveScreen) clampCursor() {
	if s.cursorRow < 0 {
		s.cursorRow = 0
	}
	if s.cursorRow >= s.rows {
		s.cursorRow = s.rows - 1
	}
	if s.cursorCol < 0 {
		s.cursorCol = 0
	}
	if s.cursorCol >= s.cols {
		s.cursorCol = s.cols - 1
	}
}

func (s *interactiveScreen) scrollUp() {
	buf := s.currentBuffer()
	copy(buf.grid[0:], buf.grid[1:])
	buf.grid[s.rows-1] = make([]interactiveCell, s.cols)
	for c := 0; c < s.cols; c++ {
		buf.grid[s.rows-1][c] = blankCell()
	}
}

func (s *interactiveScreen) newLine() {
	s.cursorRow++
	if s.cursorRow >= s.rows {
		s.scrollUp()
		s.cursorRow = s.rows - 1
	}
}

func (s *interactiveScreen) putRune(r rune) {
	if s.cursorRow < 0 || s.cursorRow >= s.rows || s.cursorCol < 0 || s.cursorCol >= s.cols {
		s.clampCursor()
	}
	buf := s.currentBuffer()
	cell := interactiveCell{
		Ch:      r,
		FG:      s.fg,
		BG:      s.bg,
		Bold:    s.bold,
		Inverse: s.inverse,
	}
	buf.grid[s.cursorRow][s.cursorCol] = cell
	s.cursorCol++
	if s.cursorCol >= s.cols {
		s.cursorCol = 0
		s.newLine()
	}
}

func (s *interactiveScreen) clearLine(mode int) {
	buf := s.currentBuffer()
	switch mode {
	case 0:
		for c := s.cursorCol; c < s.cols; c++ {
			buf.grid[s.cursorRow][c] = blankCell()
		}
	case 1:
		for c := 0; c <= s.cursorCol && c < s.cols; c++ {
			buf.grid[s.cursorRow][c] = blankCell()
		}
	case 2:
		for c := 0; c < s.cols; c++ {
			buf.grid[s.cursorRow][c] = blankCell()
		}
	}
}

func (s *interactiveScreen) clearScreen(mode int) {
	buf := s.currentBuffer()
	switch mode {
	case 0:
		s.clearLine(0)
		for r := s.cursorRow + 1; r < s.rows; r++ {
			for c := 0; c < s.cols; c++ {
				buf.grid[r][c] = blankCell()
			}
		}
	case 1:
		s.clearLine(1)
		for r := 0; r < s.cursorRow; r++ {
			for c := 0; c < s.cols; c++ {
				buf.grid[r][c] = blankCell()
			}
		}
	case 2:
		for r := 0; r < s.rows; r++ {
			for c := 0; c < s.cols; c++ {
				buf.grid[r][c] = blankCell()
			}
		}
		s.cursorRow = 0
		s.cursorCol = 0
	}
}

func (s *interactiveScreen) setAltScreen(enabled bool) {
	if s.altScreen == enabled {
		return
	}
	s.altScreen = enabled
	if enabled {
		s.altBuf = newInteractiveBuffer(s.rows, s.cols)
		s.cursorRow = 0
		s.cursorCol = 0
	} else {
		s.cursorRow = 0
		s.cursorCol = 0
	}
}

func parseCSIParams(raw string) []int {
	if raw == "" {
		return []int{0}
	}
	parts := strings.Split(raw, ";")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			out = append(out, 0)
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			out = append(out, 0)
			continue
		}
		out = append(out, n)
	}
	return out
}

func (s *interactiveScreen) applySGR(params []int) {
	if len(params) == 0 {
		params = []int{0}
	}
	for _, p := range params {
		switch {
		case p == 0:
			s.resetAttrs()
		case p == 1:
			s.bold = true
		case p == 7:
			s.inverse = true
		case p == 22:
			s.bold = false
		case p == 27:
			s.inverse = false
		case p >= 30 && p <= 37:
			s.fg = p - 30
		case p == 39:
			s.fg = 7
		case p >= 40 && p <= 47:
			s.bg = p - 40
		case p == 49:
			s.bg = 0
		case p >= 90 && p <= 97:
			s.fg = 8 + (p - 90)
		case p >= 100 && p <= 107:
			s.bg = 8 + (p - 100)
		}
	}
}

func (s *interactiveScreen) applyCSI(final byte) {
	raw := s.csiParams.String()
	params := raw
	if s.csiPrivate && strings.HasPrefix(params, "?") {
		params = strings.TrimPrefix(params, "?")
	}
	values := parseCSIParams(params)

	switch final {
	case 'A':
		n := values[0]
		if n == 0 {
			n = 1
		}
		s.cursorRow -= n
	case 'B':
		n := values[0]
		if n == 0 {
			n = 1
		}
		s.cursorRow += n
	case 'C':
		n := values[0]
		if n == 0 {
			n = 1
		}
		s.cursorCol += n
	case 'D':
		n := values[0]
		if n == 0 {
			n = 1
		}
		s.cursorCol -= n
	case 'H', 'f':
		row, col := 1, 1
		if len(values) >= 1 && values[0] != 0 {
			row = values[0]
		}
		if len(values) >= 2 && values[1] != 0 {
			col = values[1]
		}
		s.cursorRow = row - 1
		s.cursorCol = col - 1
	case 'J':
		s.clearScreen(values[0])
	case 'K':
		s.clearLine(values[0])
	case 'm':
		s.applySGR(values)
	case 's':
		s.savedRow = s.cursorRow
		s.savedCol = s.cursorCol
	case 'u':
		s.cursorRow = s.savedRow
		s.cursorCol = s.savedCol
	case 'h':
		if s.csiPrivate && params == "1049" {
			s.setAltScreen(true)
		}
	case 'l':
		if s.csiPrivate && params == "1049" {
			s.setAltScreen(false)
		}
	}
	s.clampCursor()
}

func (s *interactiveScreen) Apply(data []byte) {
	for i := 0; i < len(data); {
		b := data[i]
		switch s.escState {
		case "normal":
			switch b {
			case 0x1b:
				s.escState = "esc"
				i++
			case '\r':
				s.cursorCol = 0
				i++
			case '\n':
				s.newLine()
				i++
			case '\b':
				if s.cursorCol > 0 {
					s.cursorCol--
				}
				i++
			case '\t':
				next := ((s.cursorCol / 8) + 1) * 8
				for s.cursorCol < next && s.cursorCol < s.cols {
					s.putRune(' ')
				}
				i++
			default:
				if b < 0x20 || b == 0x7f {
					i++
					continue
				}
				r, size := utf8.DecodeRune(data[i:])
				if r == utf8.RuneError && size == 1 {
					i++
					continue
				}
				s.putRune(r)
				i += size
			}
		case "esc":
			switch b {
			case '[':
				s.escState = "csi"
				s.csiParams.Reset()
				s.csiPrivate = false
			case ']':
				s.escState = "osc"
			case '7':
				s.savedRow = s.cursorRow
				s.savedCol = s.cursorCol
				s.escState = "normal"
			case '8':
				s.cursorRow = s.savedRow
				s.cursorCol = s.savedCol
				s.clampCursor()
				s.escState = "normal"
			case 'c':
				s.reset()
				s.escState = "normal"
			default:
				s.escState = "normal"
			}
			i++
		case "csi":
			if b == '?' && s.csiParams.Len() == 0 {
				s.csiPrivate = true
				s.csiParams.WriteByte(b)
				i++
				continue
			}
			if b >= 0x40 && b <= 0x7e {
				s.applyCSI(b)
				s.escState = "normal"
				i++
				continue
			}
			s.csiParams.WriteByte(b)
			i++
		case "osc":
			if b == 0x07 {
				s.escState = "normal"
			}
			i++
		default:
			s.escState = "normal"
		}
	}
}

func (s *interactiveScreen) resize(rows, cols int) {
	if rows <= 0 {
		rows = defaultInteractiveRows
	}
	if cols <= 0 {
		cols = defaultInteractiveCols
	}
	if rows == s.rows && cols == s.cols {
		return
	}
	s.rows = rows
	s.cols = cols
	s.mainBuf = resizeInteractiveBuffer(s.mainBuf, rows, cols)
	s.altBuf = resizeInteractiveBuffer(s.altBuf, rows, cols)
	s.clampCursor()
}

func resizeInteractiveBuffer(old *interactiveBuffer, rows, cols int) *interactiveBuffer {
	buf := newInteractiveBuffer(rows, cols)
	if old == nil {
		return buf
	}
	minRows := old.rows
	if rows < minRows {
		minRows = rows
	}
	minCols := old.cols
	if cols < minCols {
		minCols = cols
	}
	for r := 0; r < minRows; r++ {
		copy(buf.grid[r][:minCols], old.grid[r][:minCols])
	}
	return buf
}

func (s *interactiveScreen) renderText() string {
	buf := s.currentBuffer()
	var out strings.Builder
	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			ch := buf.grid[r][c].Ch
			if ch == 0 {
				ch = ' '
			}
			out.WriteRune(ch)
		}
		out.WriteByte('\n')
	}
	return out.String()
}

func (s *interactiveScreen) renderMeta(generation uint64) interactiveMeta {
	return interactiveMeta{
		Rows:       s.rows,
		Cols:       s.cols,
		CursorRow:  s.cursorRow,
		CursorCol:  s.cursorCol,
		AltScreen:  s.altScreen,
		Generation: generation,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func baseColor(idx int) color.RGBA {
	palette := []color.RGBA{
		{0x00, 0x00, 0x00, 0xff},
		{0xcd, 0x00, 0x00, 0xff},
		{0x00, 0xcd, 0x00, 0xff},
		{0xcd, 0xcd, 0x00, 0xff},
		{0x00, 0x00, 0xee, 0xff},
		{0xcd, 0x00, 0xcd, 0xff},
		{0x00, 0xcd, 0xcd, 0xff},
		{0xe5, 0xe5, 0xe5, 0xff},
		{0x7f, 0x7f, 0x7f, 0xff},
		{0xff, 0x00, 0x00, 0xff},
		{0x00, 0xff, 0x00, 0xff},
		{0xff, 0xff, 0x00, 0xff},
		{0x5c, 0x5c, 0xff, 0xff},
		{0xff, 0x00, 0xff, 0xff},
		{0x00, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0xff, 0xff},
	}
	if idx < 0 || idx >= len(palette) {
		return palette[7]
	}
	return palette[idx]
}

func fillRect(img *image.RGBA, rect image.Rectangle, c color.RGBA) {
	draw.Draw(img, rect, &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func (s *interactiveScreen) renderPNG() ([]byte, error) {
	cellW := 7
	cellH := 13
	img := image.NewRGBA(image.Rect(0, 0, s.cols*cellW, s.rows*cellH))
	fillRect(img, img.Bounds(), baseColor(0))

	buf := s.currentBuffer()
	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			cell := buf.grid[r][c]
			if cell.Ch == 0 {
				cell.Ch = ' '
			}
			fg := baseColor(cell.FG)
			bg := baseColor(cell.BG)
			if cell.Inverse {
				fg, bg = bg, fg
			}
			if r == s.cursorRow && c == s.cursorCol {
				fg, bg = bg, fg
			}
			rect := image.Rect(c*cellW, r*cellH, (c+1)*cellW, (r+1)*cellH)
			fillRect(img, rect, bg)
			if cell.Ch != ' ' {
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(fg),
					Face: basicfont.Face7x13,
					Dot:  fixed.P(c*cellW, (r+1)*cellH-2),
				}
				d.DrawString(string(cell.Ch))
			}
		}
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func atomicWriteFile(path string, data []byte, mode os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func parseWinsize(raw string) (cols int, rows int, ok bool) {
	raw = strings.TrimSpace(raw)
	parts := strings.Split(strings.ToLower(raw), "x")
	if len(parts) != 2 {
		return 0, 0, false
	}
	cols, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	rows, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || cols <= 0 || rows <= 0 {
		return 0, 0, false
	}
	return cols, rows, true
}

func formatWinsize(cols, rows int) string {
	return fmt.Sprintf("%dx%d\n", cols, rows)
}

func expandInputMacros(raw string) []byte {
	type macro struct {
		name string
		data []byte
	}
	macros := []macro{
		{"enter", []byte{'\r'}},
		{"tab", []byte{'\t'}},
		{"esc", []byte{0x1b}},
		{"up", []byte{0x1b, '[', 'A'}},
		{"down", []byte{0x1b, '[', 'B'}},
		{"right", []byte{0x1b, '[', 'C'}},
		{"left", []byte{0x1b, '[', 'D'}},
		{"backspace", []byte{0x7f}},
		{"ctrl-c", []byte{0x03}},
		{"ctrl-d", []byte{0x04}},
	}
	lookup := make(map[string][]byte, len(macros))
	for _, m := range macros {
		lookup[m.name] = m.data
	}

	var out []byte
	for i := 0; i < len(raw); {
		switch raw[i] {
		case '\\':
			if i+1 < len(raw) {
				out = append(out, raw[i+1])
				i += 2
				continue
			}
			out = append(out, raw[i])
			i++
		case '<':
			j := strings.IndexByte(raw[i:], '>')
			if j <= 0 {
				out = append(out, raw[i])
				i++
				continue
			}
			token := strings.ToLower(raw[i+1 : i+j])
			if seq, ok := lookup[token]; ok {
				out = append(out, seq...)
			} else {
				out = append(out, raw[i:i+j+1]...)
			}
			i += j + 1
		default:
			out = append(out, raw[i])
			i++
		}
	}
	return out
}

func (b *ShExecutor) renderInteractiveSnapshot(state *interactiveState) {
	state.mu.Lock()
	state.generation++
	screenText := state.screen.renderText()
	meta := state.screen.renderMeta(state.generation)
	pngBytes, pngErr := state.screen.renderPNG()
	state.mu.Unlock()

	_ = atomicWriteFile(state.screenPath, []byte(screenText), 0o644)
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	_ = atomicWriteFile(state.metaPath, append(metaBytes, '\n'), 0o644)
	if pngErr == nil {
		_ = atomicWriteFile(state.pngPath, pngBytes, 0o644)
	}
}

func (b *ShExecutor) pollInteractiveInput(job *managedJob) {
	state := job.ptyState
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	defer state.wg.Done()

	for {
		select {
		case <-state.stopCh:
			return
		case <-ticker.C:
			data, err := os.ReadFile(state.inputPath)
			if err != nil || len(data) == 0 {
				continue
			}
			expanded := expandInputMacros(string(data))
			if len(expanded) > 0 {
				if state.inputLog != nil {
					_, _ = state.inputLog.Write(expanded)
				}
				_, _ = state.ptyFile.Write(expanded)
			}
			_ = os.WriteFile(state.inputPath, nil, 0o644)
		}
	}
}

func (b *ShExecutor) pollInteractiveWinsize(job *managedJob) {
	state := job.ptyState
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer state.wg.Done()

	for {
		select {
		case <-state.stopCh:
			return
		case <-ticker.C:
			data, err := os.ReadFile(state.winsizePath)
			if err != nil {
				continue
			}
			raw := strings.TrimSpace(string(data))
			if raw == "" || raw == state.lastWinsize {
				continue
			}
			cols, rows, ok := parseWinsize(raw)
			if !ok {
				continue
			}
			ws := &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)}
			if err := pty.Setsize(state.ptyFile, ws); err != nil {
				continue
			}
			state.mu.Lock()
			state.screen.resize(rows, cols)
			state.lastWinsize = strings.TrimSpace(formatWinsize(cols, rows))
			state.mu.Unlock()
			_ = os.WriteFile(state.winsizePath, []byte(formatWinsize(cols, rows)), 0o644)
			b.renderInteractiveSnapshot(state)
		}
	}
}

func (b *ShExecutor) readInteractivePTY(job *managedJob) {
	state := job.ptyState
	defer state.wg.Done()

	buf := make([]byte, 4096)
	for {
		n, err := state.ptyFile.Read(buf)
		if n > 0 {
			chunk := append([]byte(nil), buf[:n]...)
			if state.eventsFile != nil {
				_, _ = state.eventsFile.Write(chunk)
			}
			if state.eventsHex != nil {
				_, _ = state.eventsHex.WriteString(hex.EncodeToString(chunk) + "\n")
			}
			state.mu.Lock()
			state.screen.Apply(chunk)
			state.mu.Unlock()
			b.renderInteractiveSnapshot(state)
		}
		if err != nil {
			if err != io.EOF {
				// ignore short-lived PTY read errors once the process exits
			}
			return
		}
	}
}

func (b *ShExecutor) awaitInteractiveExit(job *managedJob) {
	state := job.ptyState
	code := exitCodeFromWait(job.cmd.Wait())
	job.exitCode = code

	close(state.stopCh)
	_ = state.ptyFile.Close()
	state.wg.Wait()
	b.renderInteractiveSnapshot(state)

	b.finalizeInteractiveWorkspace(job)
	status := fmt.Sprintf("%d\n", code)
	if err := atomicWriteFile(job.exitPath, []byte(status), 0o644); err != nil {
		publishExitError(job, fmt.Sprintf("failed to publish exit code %d: %v", code, err))
	}
	if state.eventsFile != nil {
		_ = state.eventsFile.Close()
	}
	if state.eventsHex != nil {
		_ = state.eventsHex.Close()
	}
	if state.inputLog != nil {
		_ = state.inputLog.Close()
	}

	close(job.doneCh)
	if job.detached {
		b.mu.Lock()
		delete(b.detached, job.ID)
		b.mu.Unlock()
	}
}

func (b *ShExecutor) finalizeInteractiveWorkspace(job *managedJob) {
	if job == nil || job.workspace == nil || !job.workspace.enabled || !job.workspace.usesOverlayBackend() {
		return
	}
	finalized, err := job.workspace.finalizeTurn("interactive", b.TurnID)
	if err != nil {
		_ = atomicWriteFile(filepath.Join(job.canonicalDir, "workspace_error"), []byte(err.Error()+"\n"), 0o644)
		return
	}
	if b.FSMutationTelemetryEnabled {
		_ = atomicWriteFile(filepath.Join(job.canonicalDir, "fs_mutations"), []byte(finalized.Mutations+"\n"), 0o644)
	}
	if finalized.Revision.ID == "" {
		return
	}
	worldRevisionBlock := formatWorldRevisionCreated(finalized.Revision, !finalized.Changed)
	if worldRevisionBlock != "" {
		_ = atomicWriteFile(filepath.Join(job.canonicalDir, "world_revision"), []byte(worldRevisionBlock+"\n"), 0o644)
	}
	handle := buildWorldHandle(job.workspace.workspaceSession, finalized.Revision.ID)
	_ = atomicWriteFile(filepath.Join(job.canonicalDir, "world_handle"), []byte(handle+"\n"), 0o644)
}

func (b *ShExecutor) newInteractiveJobWorkspace() (*subjectiveFS, error) {
	if b.subjective == nil || !b.subjective.enabled || !b.subjective.usesOverlayBackend() {
		return nil, nil
	}
	parentSession := strings.TrimSpace(b.subjective.workspaceSession)
	if parentSession == "" {
		parentSession = strings.TrimSpace(b.SessionID)
	}
	jobSession := newInteractiveWorkspaceSessionID(parentSession)
	workspace := &subjectiveFS{
		enabled:          true,
		sessionID:        jobSession,
		workspaceRoot:    b.subjective.workspaceRoot,
		workspace:        b.subjective.workspace,
		workspaceBackend: b.subjective.workspaceBackend,
		revisionMode:     b.subjective.revisionMode,
		dataDir:          b.subjective.dataDir,
		workspaceSession: jobSession,
		workspaceOwner:   false,
		bootstrapSource:  b.subjective.workspaceSession,
	}
	if err := workspace.init(b.DataDir, jobSession); err != nil {
		return nil, err
	}
	return workspace, nil
}

func newInteractiveWorkspaceSessionID(parent string) string {
	parent = strings.TrimSpace(parent)
	if parent == "" {
		parent = "workspace"
	}
	parent = strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(parent)
	return fmt.Sprintf("%s-interactive-%d", parent, time.Now().UnixNano())
}

func (b *ShExecutor) startInteractiveJob(command string) (*managedJob, error) {
	if err := b.initState(); err != nil {
		return nil, err
	}

	jobSessionRoot, err := b.jobSessionRoot()
	if err != nil {
		return nil, err
	}
	jobWorkspace, err := b.newInteractiveJobWorkspace()
	if err != nil {
		return nil, fmt.Errorf("creating interactive workspace: %w", err)
	}

	cmd := exec.Command(b.Shell, "-c", jobWrapperScript)
	cmd.Env = MergeEnv(b.commandBaseEnv(), jobWrapperEnv(b.Shell, command, jobSessionRoot, b.Network,
		"QUINE_JOB_INTERACTIVE=1",
	))
	useOverlayWorkspace := jobWorkspace != nil && jobWorkspace.enabled && jobWorkspace.usesOverlayBackend()
	isolateNetwork := b.Network == "none"
	if useOverlayWorkspace {
		cmd.Dir = jobSessionRoot
		cmd.Env = MergeEnv(cmd.Env, jobWorkspace.commandEnv())
	} else if b.WorkDir != "" {
		cmd.Dir = b.WorkDir
	} else if b.subjective != nil && b.subjective.enabled && b.subjective.usesDirectBackend() && strings.TrimSpace(b.subjective.workspace) != "" {
		cmd.Dir = b.subjective.workspace
	} else {
		cmd.Dir = jobSessionRoot
	}
	cmd.ExtraFiles = b.extraFiles()

	ws := &pty.Winsize{Cols: defaultInteractiveCols, Rows: defaultInteractiveRows}
	ptyFile, err := pty.StartWithAttrs(cmd, ws, jobSysProcAttr(true, useOverlayWorkspace, isolateNetwork))
	if err != nil {
		return nil, fmt.Errorf("starting interactive process: %w", err)
	}

	pid := cmd.Process.Pid
	startedAt := time.Now().UTC()
	cleanupStartFailure := func() {
		_ = ptyFile.Close()
		_ = syscall.Kill(-pid, syscall.SIGKILL)
		_ = cmd.Wait()
	}
	canonicalDir, err := stageJobSurface(jobSessionRoot, pid, func(dir string) error {
		if err := writeJobIdentityFiles(dir, command, pid, startedAt); err != nil {
			return err
		}
		for _, p := range []string{
			filepath.Join(dir, "in"),
			filepath.Join(dir, "screen.txt"),
			filepath.Join(dir, "screen.meta"),
			filepath.Join(dir, "screen.png"),
			filepath.Join(dir, "winsize"),
			filepath.Join(dir, "events.log"),
			filepath.Join(dir, "events.hex"),
			filepath.Join(dir, "input.log"),
		} {
			if err := touchFile(p); err != nil {
				return fmt.Errorf("initializing interactive file %s: %w", p, err)
			}
		}
		if jobWorkspace != nil {
			if err := os.WriteFile(filepath.Join(dir, "workspace_session"), []byte(jobWorkspace.workspaceSession+"\n"), 0o644); err != nil {
				return fmt.Errorf("writing workspace_session file: %w", err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, "winsize"), []byte(formatWinsize(defaultInteractiveCols, defaultInteractiveRows)), 0o644); err != nil {
			return fmt.Errorf("writing winsize file: %w", err)
		}
		return nil
	})
	if err != nil {
		cleanupStartFailure()
		return nil, err
	}

	inputPath := filepath.Join(canonicalDir, "in")
	screenPath := filepath.Join(canonicalDir, "screen.txt")
	metaPath := filepath.Join(canonicalDir, "screen.meta")
	pngPath := filepath.Join(canonicalDir, "screen.png")
	winsizePath := filepath.Join(canonicalDir, "winsize")
	eventsPath := filepath.Join(canonicalDir, "events.log")
	eventsHexPath := filepath.Join(canonicalDir, "events.hex")
	inputLogPath := filepath.Join(canonicalDir, "input.log")
	exitPath := filepath.Join(canonicalDir, "exit")

	eventsFile, err := os.OpenFile(eventsPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		cleanupStartFailure()
		_ = os.RemoveAll(canonicalDir)
		return nil, fmt.Errorf("opening events.log: %w", err)
	}
	eventsHex, err := os.OpenFile(eventsHexPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		_ = eventsFile.Close()
		cleanupStartFailure()
		_ = os.RemoveAll(canonicalDir)
		return nil, fmt.Errorf("opening events.hex: %w", err)
	}
	inputLog, err := os.OpenFile(inputLogPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		_ = eventsHex.Close()
		_ = eventsFile.Close()
		cleanupStartFailure()
		_ = os.RemoveAll(canonicalDir)
		return nil, fmt.Errorf("opening input.log: %w", err)
	}

	state := &interactiveState{
		ptyFile:     ptyFile,
		eventsFile:  eventsFile,
		eventsHex:   eventsHex,
		inputLog:    inputLog,
		inputPath:   inputPath,
		screenPath:  screenPath,
		metaPath:    metaPath,
		pngPath:     pngPath,
		winsizePath: winsizePath,
		eventsPath:  eventsPath,
		screen:      newInteractiveScreen(defaultInteractiveRows, defaultInteractiveCols),
		lastWinsize: strings.TrimSpace(formatWinsize(defaultInteractiveCols, defaultInteractiveRows)),
		stopCh:      make(chan struct{}),
	}
	job := &managedJob{
		ID:           pid,
		cmd:          cmd,
		detached:     true,
		interactive:  true,
		canonicalDir: canonicalDir,
		displayDir:   filepath.ToSlash(canonicalDir) + "/",
		exitPath:     exitPath,
		doneCh:       make(chan struct{}),
		ptyState:     state,
		workspace:    jobWorkspace,
	}

	b.renderInteractiveSnapshot(state)

	state.wg.Add(3)
	go b.readInteractivePTY(job)
	go b.pollInteractiveInput(job)
	go b.pollInteractiveWinsize(job)
	go b.awaitInteractiveExit(job)

	return job, nil
}
