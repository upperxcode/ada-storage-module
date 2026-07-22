package storage

import "time"

// RawMessage represents a complete message from the AI provider,
// including thinking content and tool calls.
// This type is used for integration with ada-love-core.
type RawMessage struct {
	Role             string    `json:"role"`
	Content          string    `json:"content"`
	ToolCalls        []any     `json:"tool_calls"`
	ToolCallID       string    `json:"tool_call_id"`
	Time             time.Time `json:"time"`
	ThinkingContent  string    `json:"thinking_content,omitempty"`
	ThinkingDuration int       `json:"thinking_duration,omitempty"`
}
