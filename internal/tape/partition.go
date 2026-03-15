package tape

import "encoding/json"

// EntryPartition is the logical plane a tape entry belongs to.
type EntryPartition string

const (
	PartitionMeta    EntryPartition = "meta"
	PartitionSystem  EntryPartition = "system"
	PartitionContext EntryPartition = "context"
)

// PartitionOf classifies a tape entry into meta, system, or context.
func PartitionOf(entry TapeEntry) EntryPartition {
	switch entry.Type {
	case "meta":
		return PartitionMeta
	case "message":
		var msg Message
		if err := json.Unmarshal(entry.Data, &msg); err == nil && msg.Role == RoleSystem {
			return PartitionSystem
		}
	}
	return PartitionContext
}

// ContextEntries returns only entries that belong to the context plane.
func ContextEntries(entries []TapeEntry) []TapeEntry {
	out := make([]TapeEntry, 0, len(entries))
	for _, entry := range entries {
		if PartitionOf(entry) == PartitionContext {
			out = append(out, entry)
		}
	}
	return out
}
