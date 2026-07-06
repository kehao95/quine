package qcli

type HelloEvent struct {
	Type       string   `json:"type"`
	Contract   string   `json:"contract"`
	EndpointID string   `json:"endpoint_id"`
	Peer       PeerInfo `json:"peer"`
}

type StatusEvent struct {
	Type          string `json:"type,omitempty"`
	Live          bool   `json:"live"`
	Working       bool   `json:"working"`
	PID           int    `json:"pid"`
	PPID          int    `json:"ppid"`
	Session       string `json:"session"`
	RunID         string `json:"run_id"`
	Incarnation   int    `json:"incarnation"`
	ParentSession string `json:"parent_session"`
	Model         string `json:"model"`
	Depth         int    `json:"depth"`
	Pending       int    `json:"pending"`
	TokensIn      *int   `json:"tokens_in"`
	TokensOut     *int   `json:"tokens_out"`
}

type Cell struct {
	Type            string  `json:"type,omitempty"`
	Seq             int64   `json:"seq"`
	Kind            string  `json:"kind"`
	Role            *string `json:"role"`
	Text            *string `json:"text"`
	TS              *int64  `json:"ts"`
	ToolID          *string `json:"tool_id"`
	ToolName        *string `json:"tool_name"`
	Status          *string `json:"status"`
	IsError         bool    `json:"is_error"`
	Action          *string `json:"action"`
	Author          *string `json:"author"`
	Raw             *string `json:"raw"`
	Depth           *int    `json:"depth,omitempty"`
	Session         *string `json:"session,omitempty"`
	TerminationMode *string `json:"termination_mode,omitempty"`
	ExitCode        *int    `json:"exit_code,omitempty"`
	TurnCount       *int    `json:"turn_count,omitempty"`
	TokensIn        *int    `json:"tokens_in,omitempty"`
	TokensOut       *int    `json:"tokens_out,omitempty"`
}

type BackfillCompleteEvent struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type StreamDeltaEvent struct {
	Type       string  `json:"type"`
	Generation int     `json:"generation"`
	Seq        int     `json:"seq"`
	Kind       string  `json:"kind"`
	Text       *string `json:"text"`
	ToolID     *string `json:"tool_id"`
	ToolName   *string `json:"tool_name"`
	TS         int64   `json:"ts"`
}

type ReceiptEvent struct {
	Type      string  `json:"type"`
	Stage     string  `json:"stage"`
	Action    *string `json:"action"`
	Delivery  *string `json:"delivery"`
	MessageID *string `json:"message_id"`
	ClientRef *string `json:"client_ref"`
	Pending   *int    `json:"pending"`
	TS        int64   `json:"ts"`
}

type ContextResetEvent struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	TS     int64  `json:"ts"`
}

type ErrorEvent struct {
	Type        string `json:"type"`
	Scope       string `json:"scope"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	Recoverable bool   `json:"recoverable"`
}

type QueuedResponse struct {
	Type      string `json:"type"`
	Action    string `json:"action"`
	ClientRef string `json:"client_ref"`
}

type AttachedResponse struct {
	Type string   `json:"type"`
	Peer PeerInfo `json:"peer"`
}

type RosterResponse struct {
	Type  string        `json:"type"`
	Peers []RosterEntry `json:"peers"`
}

type PeerContractResponse struct {
	Type     string `json:"type"`
	Contract any    `json:"contract"`
}

func strptr(s string) *string {
	return &s
}

func intptr(v int) *int {
	return &v
}

func int64ptr(v int64) *int64 {
	return &v
}

func withCellType(c Cell) Cell {
	c.Type = "cell"
	return c
}
