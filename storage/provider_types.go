package storage

import (
	"database/sql"
	"errors"
)

// Config-related errors
var (
	ErrConfigNotFound = errors.New("configuration key not found")
)

// Provider-related errors
var (
	ErrProviderNotFound = errors.New("provider not found")
	ErrProviderModelNotFound = errors.New("provider model not found")
	ErrGreetingNotFound = errors.New("greeting not found")
)

// Provider represents an LLM provider.
type Provider struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	APIURL         sql.NullString `json:"api_url"`
	ConnectionTypes sql.NullString `json:"connection_types"`
	Color          string         `json:"color"`
	Icon           string         `json:"icon"`
}

// ProviderModel represents a model available from a provider.
type ProviderModel struct {
	ID        int64  `json:"id"`
	ProviderID int64 `json:"provider_id"`
	Model     string `json:"model"`
	Free      bool   `json:"free"`
	Thinking  bool   `json:"thinking"`
	Tool      bool   `json:"tool"`
	Embedding bool   `json:"embedding"`
	Vision    bool   `json:"vision"`
	Health    int    `json:"health"`
}

// Greeting represents a greeting response for a keyword.
type Greeting struct {
	ID       int64  `json:"id"`
	Keyword  string `json:"keyword"`
	Language string `json:"language"`
	Response string `json:"response"`
}

// Config represents a configuration key-value pair.
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}