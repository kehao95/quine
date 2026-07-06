package qcli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ReadPeerContract(agent Agent) (any, error) {
	path := filepath.Join(agent.PublicRoot, "status", "contract.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("qcli: parse peer contract: %w", err)
	}
	return v, nil
}
