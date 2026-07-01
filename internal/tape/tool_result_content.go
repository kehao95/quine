package tape

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MarshalToolResultContent encodes a tool-result payload as a compact JSON
// object. On marshal failure, it returns a compact JSON error object rather
// than falling back to non-JSON content.
func MarshalToolResultContent(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err == nil && json.Valid(data) {
		var compact bytes.Buffer
		if compactErr := json.Compact(&compact, data); compactErr == nil {
			return json.RawMessage(compact.Bytes())
		}
		return json.RawMessage(data)
	}

	fallback, fallbackErr := json.Marshal(map[string]any{
		"tool":   "runtime",
		"status": "error",
		"error":  fmt.Sprintf("failed to marshal tool result content: %v", err),
	})
	if fallbackErr != nil {
		return json.RawMessage(`{"tool":"runtime","status":"error","error":"failed to marshal tool result content"}`)
	}
	return json.RawMessage(fallback)
}

func decodeToolResultPayload(data []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var payload any
	if err := dec.Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode tool result content: %w", err)
	}

	if dec.More() {
		return nil, fmt.Errorf("tool result content must be a single JSON object")
	}

	obj, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("tool result content must be a JSON object")
	}
	return obj, nil
}

// UnmarshalToolResultContent decodes a tool-result content string and enforces
// that the payload is a JSON object.
func UnmarshalToolResultContent(content string) (map[string]any, error) {
	return decodeToolResultPayload([]byte(content))
}

// UnmarshalStructuredToolResultContent decodes persisted object-valued tool
// results and enforces that the payload is a JSON object.
func UnmarshalStructuredToolResultContent(content json.RawMessage) (map[string]any, error) {
	return decodeToolResultPayload(content)
}

// RewriteToolResultContent decodes, mutates, and re-encodes tool-result
// content while preserving the single-object JSON contract.
func RewriteToolResultContent(content string, mutate func(map[string]any) error) (string, error) {
	obj, err := UnmarshalToolResultContent(content)
	if err != nil {
		return "", err
	}
	if err := mutate(obj); err != nil {
		return "", err
	}
	return string(MarshalToolResultContent(obj)), nil
}

// RewriteStructuredToolResultContent decodes, mutates, and re-encodes
// persisted object-valued tool results while preserving the single-object JSON
// contract.
func RewriteStructuredToolResultContent(content json.RawMessage, mutate func(map[string]any) error) (json.RawMessage, error) {
	obj, err := UnmarshalStructuredToolResultContent(content)
	if err != nil {
		return nil, err
	}
	if err := mutate(obj); err != nil {
		return nil, err
	}
	return MarshalToolResultContent(obj), nil
}
