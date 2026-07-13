# ada-storage-module

An independent, highly decoupled SQLite storage package for the ada-love-ai backend.

## Philosophy

This module follows strict Fail-Fast principles:
- Errors are never masked or swallowed
- Migration failures abort application startup
- Missing data returns explicit errors (`ErrSessionNotFound`, `ErrConfigNotFound`, etc.)

## Architecture

The package is organized by single responsibility, with all files under 300 lines:

```
storage/
├── database.go               # Connection initialization & pool management
├── init.go                   # Module initialization entry point
├── migrations.go             # Migration runner with version tracking
├── migrations_core.go        # Core tables (sessions, messages, app_state, router_configs)
├── migrations_agents.go      # Agent-related tables (agents, workers, skills)
├── migrations_config.go      # Config tables (config, greetings)
├── migrations_providers.go   # Provider tables (providers, models, apikeys)
├── migrations_workspaces.go  # Workspace tables (workspaces, memories, associations)
├── migrations_fixed.go       # Fixed models and tool profiles
├── migrations_other.go       # Other tables (mcps, spec_wizards, workspace_templates)
├── session_store.go          # Sessions and messages operations
├── workspace_store.go        # Workspace operations
├── memory_store.go           # Memory and workspace-agent operations
├── agent_types.go            # Agent, Worker, Skill types
├── provider_types.go         # Provider, ProviderModel, Greeting, Config types
├── agent_store.go            # Agent operations
├── worker_store.go           # Worker operations
├── skill_store.go            # Skill operations
├── provider_store.go         # Provider operations
├── provider_model_store.go   # Provider model operations
├── greeting_store.go         # Greeting operations
└── config_store.go           # Generic config and API key operations
```

## Usage

### Initialization

```go
package main

import (
    "context"
    "log"
    
    "github.com/user/ada-storage-module/storage"
)

func main() {
    ctx := context.Background()
    
    // Initialize storage with default SQLite database
    engine, sessionStore, configStore, err := storage.Init(ctx, "./data/ada.db")
    if err != nil {
        log.Fatalf("Failed to initialize storage: %v", err)
    }
    defer engine.Close()
    
    // Use stores...
    _ = sessionStore
    _ = configStore
}
```

### Session Operations

```go
// Create a new session
sessionID := storage.GenerateSessionID()
err := sessionStore.CreateSession(ctx, &storage.Session{
    ID:    sessionID,
    Title: sql.NullString{String: "My Session", Valid: true},
})

// Add a message
err = sessionStore.SaveMessage(ctx, &storage.Message{
    SessionID: sessionID,
    Role:      "user",
    Content:   "Hello, world!",
})

// Retrieve session history
messages, err := sessionStore.GetMessages(ctx, sessionID)
```

### Provider Operations

```go
// Create a provider
provider := &storage.Provider{Name: "openai"}
err := providerStore.CreateProvider(ctx, provider)

// Get provider models
models, err := providerModelStore.GetProviderModels(ctx, provider.ID)

// Set API key
err = configStore.SetAPIKey(ctx, provider.ID, "sk-your-key")
```

### Config Operations

```go
// Set/Get configuration
err := configStore.SetConfig(ctx, "api_key", "secret-value")
value, err := configStore.GetConfig(ctx, "api_key")

// Load all greetings
greetings, err := greetingStore.LoadAllGreetings(ctx)
```

### Workspace Operations

```go
// Create a workspace
workspace := &storage.Workspace{
    Nome:  "My Project",
    Path:  sql.NullString{String: "/path/to/workspace", Valid: true},
    Enabled: true,
}
err := workspaceStore.CreateWorkspace(ctx, workspace)

// Add memory
memory := &storage.Memory{
    Content: "Important context",
    Embedding: []byte{...},
}
err := memoryStore.CreateMemory(ctx, memory)
```

## Error Handling

All errors are explicit and never masked:

```go
// Session not found error
_, err := sessionStore.GetSession(ctx, "nonexistent-id")
// err == storage.ErrSessionNotFound

// Provider not found error
_, err := providerStore.GetProvider(ctx, 999)
// err == storage.ErrProviderNotFound

// Config key not found error
_, err := configStore.GetConfig(ctx, "missing-key")
// err == storage.ErrConfigNotFound
```

## Database Schema

### Core Tables
- **sessions** - Chat sessions with metadata
- **messages** - Chat messages with foreign key to sessions
- **app_state** - Application state (singleton)
- **router_configs** - Router configuration

### Provider Tables
- **providers** - LLM providers
- **provider_models** - Models per provider
- **provider_apikeys** - API keys per provider

### Config Tables
- **config** - Generic key-value configuration
- **greetings** - Keyword-triggered greetings

### Workspace Tables
- **workspaces** - Workspace metadata
- **memories** - Workspace memories with embeddings
- **workspace_agents** - Workspace-agent associations
- **workspace_folders** - Workspace folder structure
- **workspace_skills** - Workspace-skill associations
- **workspace_tools** - Workspace tool configurations
- **workspace_workers** - Workspace-worker associations

### Agent Tables
- **agents** - Agent definitions
- **workers** - Worker definitions
- **skills** - Skill definitions

## Dependencies

- `github.com/mattn/go-sqlite3` - SQLite driver with CGO
- `github.com/google/uuid` - UUID generation

## Testing

```bash
go test ./storage/... -v
```

## License

MIT