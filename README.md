# ada-storage-module

A robust, independent SQLite storage package for the ada-love-ai backend.

## 🎯 Objetivo

This module provides a highly decoupled storage layer for managing:
- Chat sessions and messages
- Provider models and API keys
- System configurations and greetings
- Workspaces with memories
- Agents, workers, and skills

### ✨ Princípios

1. **Fail-Fast**: Errors are never masked; migration failures abort startup
2. **Single Responsibility**: Each file handles one domain (under 300 lines)
3. **Explicit Errors**: Custom error types (`ErrSessionNotFound`, `ErrProviderNotFound`, etc.)
4. **WAL Mode**: Write-Ahead Logging for concurrent read/write performance
5. **Foreign Keys**: Referential integrity with CASCADE deletes

---

## 📦 Arquitetura

```
storage/
├── database.go               # Connection pool + WAL/foreign keys
├── init.go                   # Module initialization entry point
├── migrations.go             # Migration runner (v1-v29)
├── migrations_*.go           # Schema definitions by category
├── session_store.go          # Sessions & messages
├── workspace_store.go        # Workspaces
├── memory_store.go           # Memories & workspace agents
├── agent_types.go            # Type definitions
├── provider_types.go         # Type definitions
├── agent_store.go            # Agents
├── worker_store.go           # Workers
├── skill_store.go            # Skills
├── provider_store.go         # Providers
├── provider_model_store.go   # Provider models
├── greeting_store.go         # Greetings
└── config_store.go           # Config & API keys
```

---

## 🔧 Inicialização

```go
package main

import (
    "context"
    "log"
    
    "github.com/ada-love-ai/storage/storage"
)

func main() {
    ctx := context.Background()
    
    // Initialize storage
    engine, sessionStore, configStore, err := storage.Init(ctx, "./data/ada.db")
    if err != nil {
        log.Fatalf("Failed to initialize: %v", err)
    }
    defer engine.Close()
    
    // Use stores...
    _ = sessionStore
    _ = configStore
}
```

### 📁 Caminhos de Banco

| Valor | Uso |
|-------|-----|
| `"./data/ada.db"` | Banco local padrão |
| `":memory:"` | Banco em memória (testes) |
| `"file::memory:?cache=shared"` | Banco compartilhado em memória |
| `"/var/data/ada.db"` | Caminho absoluto |

---

## 📚 API Reference

### SessionStore - Sessões e Mensagens

#### `NewSessionStore(db *sql.DB) *SessionStore`
Cria uma nova instância do SessionStore.

#### `GenerateSessionID() string`
Gera um ID único para sessão usando UUID.

#### `CreateSession(ctx, *Session) error`
Cria uma nova sessão.

```go
err := sessionStore.CreateSession(ctx, &storage.Session{
    ID:    storage.GenerateSessionID(),
    Title: sql.NullString{String: "My Session", Valid: true},
})
```

#### `GetSession(ctx, id string) (*Session, error)`
Recupera uma sessão por ID. Retorna `ErrSessionNotFound` se não existir.

```go
session, err := sessionStore.GetSession(ctx, sessionID)
if errors.Is(err, storage.ErrSessionNotFound) {
    // Handle missing session
}
```

#### `SaveMessage(ctx, *Message) error`
Salva uma mensagem na sessão.

```go
err := sessionStore.SaveMessage(ctx, &storage.Message{
    SessionID: sessionID,
    Role:      "user",
    Content:   "Hello, world!",
})
```

#### `GetMessages(ctx, sessionID string) ([]Message, error)`
Recupera todas as mensagens de uma sessão em ordem cronológica.

```go
messages, err := sessionStore.GetMessages(ctx, sessionID)
for _, msg := range messages {
    fmt.Printf("[%s] %s\n", msg.Role, msg.Content)
}
```

#### `DeleteSession(ctx, id string) error`
Apaga uma sessão e todas as mensagens (CASCADE).

#### `ListSessions(ctx) ([]Session, error)`
Lista todas as sessões ordenadas por data de atualização.

---

### ProviderStore - Providers

#### `CreateProvider(ctx, *Provider) error`
Cria um novo provider.

```go
err := providerStore.CreateProvider(ctx, &storage.Provider{
    Name: "openai",
})
```

#### `GetProvider(ctx, id int64) (*Provider, error)`
Recupera um provider por ID.

#### `GetProviderByName(ctx, name string) (*Provider, error)`
Recupera um provider pelo nome.

#### `ListProviders(ctx) ([]Provider, error)`
Lista todos os providers.

#### `DeleteProvider(ctx, id int64) error`
Apaga um provider.

---

### ProviderModelStore - Modelos de Provider

#### `CreateProviderModel(ctx, *ProviderModel) error`
Cria um novo modelo para um provider.

```go
err := modelStore.CreateProviderModel(ctx, &storage.ProviderModel{
    ProviderID: providerID,
    Model:      "gpt-4",
    Free:       false,
    Thinking:   true,
    Tool:       true,
    Embedding:  false,
    Vision:     true,
    Health:     100,
})
```

#### `GetProviderModels(ctx, providerID int64) ([]ProviderModel, error)`
Recupera todos os modelos de um provider.

#### `GetEnabledModelsByProvider(ctx) (map[string][]ProviderModel, error)`
Recupera modelos habilitados agrupados por nome de provider.

---

### ConfigStore - Configurações

#### `SetConfig(ctx, key, value string) error`
Salva uma configuração.

```go
err := configStore.SetConfig(ctx, "api_key", "secret-value")
```

#### `GetConfig(ctx, key string) (string, error)`
Recupera uma configuração. Retorna `ErrConfigNotFound` se não existir.

```go
value, err := configStore.GetConfig(ctx, "api_key")
if errors.Is(err, storage.ErrConfigNotFound) {
    // Handle missing config
}
```

#### `DeleteConfig(ctx, key string) error`
Apaga uma configuração.

#### `ListConfigs(ctx) (map[string]string, error)`
Lista todas as configurações.

#### `SetAPIKey(ctx, providerID int64, apiKey string) error`
Salva uma API key para um provider.

#### `GetAPIKeys(ctx, providerID int64) ([]string, error)`
Recupera todas as API keys de um provider.

---

### GreetingStore - Saudações

#### `CreateGreeting(ctx, *Greeting) error`
Cria uma saudação para uma palavra-chave.

```go
err := greetingStore.CreateGreeting(ctx, &storage.Greeting{
    Keyword:  "hello",
    Language: "en",
    Response: "Hello! How can I help you?",
})
```

#### `GetGreeting(ctx, keyword, language string) (*Greeting, error)`
Recupera uma saudação específica.

#### `LoadAllGreetings(ctx) (map[string]string, error)`
Carrega todas as saudações em um mapa (usado no boot).

```go
greetings, err := greetingStore.LoadAllGreetings(ctx)
response := greetings["hello"] // Retorna a saudação em português se disponível
```

---

### WorkspaceStore - Workspaces

#### `CreateWorkspace(ctx, *Workspace) error`
Cria um novo workspace.

```go
err := workspaceStore.CreateWorkspace(ctx, &storage.Workspace{
    Nome:    "My Project",
    Path:    sql.NullString{String: "/path/to/workspace", Valid: true},
    Enabled: true,
})
```

#### `GetWorkspace(ctx, id int64) (*Workspace, error)`
Recupera um workspace por ID.

#### `GetWorkspaceByPath(ctx, path string) (*Workspace, error)`
Recupera um workspace pelo caminho.

#### `ListWorkspaces(ctx) ([]Workspace, error)`
Lista todos os workspaces.

#### `DeleteWorkspace(ctx, id int64) error`
Apaga um workspace e todos os dados associados (CASCADE).

---

### MemoryStore - Memórias

#### `CreateMemory(ctx, *Memory) error`
Cria uma nova memória.

```go
err := memoryStore.CreateMemory(ctx, &storage.Memory{
    WorkspacePath: sql.NullString{String: workspacePath, Valid: true},
    Content:       "Important context",
    Importance:    1,
    Embedding:   []byte{...},
})
```

#### `GetMemories(ctx, workspacePath string) ([]Memory, error)`
Recupera todas as memórias de um workspace.

#### `DeleteMemories(ctx, workspacePath string) error`
Apaga todas as memórias de um workspace.

---

### AgentStore - Agents

#### `CreateAgent(ctx, *Agent) error`
Cria um novo agent.

```go
err := agentStore.CreateAgent(ctx, &storage.Agent{
    Name:        "Research Agent",
    Description: sql.NullString{String: "Research assistant", Valid: true},
    Type:        storage.AgentTypeResearch,
    MaxIteration: 10,
    Temperature: 0.7,
})
```

#### `ListAgents(ctx) ([]Agent, error)`
Lista todos os agents.

---

### WorkerStore - Workers

#### `CreateWorker(ctx, *Worker) error`
Cria um novo worker.

```go
err := workerStore.CreateWorker(ctx, &storage.Worker{
    Name:             "CLI Worker",
    ConnectionType:   "cli_command",
    Command:          sql.NullString{String: "python script.py", Valid: true},
    ResponseLanguage: "portuguese",
})
```

---

### SkillStore - Skills

#### `CreateSkill(ctx, *Skill) error`
Cria uma nova skill.

```go
err := skillStore.CreateSkill(ctx, &storage.Skill{
    Name:        "search",
    Description: sql.NullString{String: "Search functionality", Valid: true},
    Content:     "def search(query): ...",
})
```

#### `SearchSkills(ctx, keyword string) ([]Skill, error)`
Pesquisa skills por palavra-chave.

---

## 🛠️ Scripts de Desenvolvimento

```bash
# Menu interativo com cores e emojis
./scripts/run

# Opções:
# 1) run          - Executar aplicacao
# 2) build        - Compilar projeto
# 3) test         - Testes unitarios
# 4) test-all     - Todos os testes
# 5) vet          - Executar go vet
# 6) fmt          - Formatar codigo
# 7) clean        - Limpar arquivos
# 8) dev          - Modo desenvolvimento
# Q) Sair
```

---

## 🧪 Testes

```bash
# Executar todos os testes
go test ./storage/... -v

# Com coverage
go test -cover ./storage/...

# Com race detection
go test -race ./storage/...
```

---

## 📊 Estatísticas

| Métrica | Valor |
|---------|-------|
| Arquivos | 25 |
| Linhas Totais | 2.600 |
| Tabelas | 34 |
| Migrações | 29 |
| Máx linhas/arquivo | 280 |

---

## 📝 Dependências

```bash
github.com/google/uuid v1.6.0
github.com/mattn/go-sqlite3 v1.14.48
```

---

## ⚠️ Erros Comuns

| Erro | Causa | Solução |
|------|-------|---------|
| `ErrSessionNotFound` | Sessão não existe | Verificar ID ou criar nova sessão |
| `ErrProviderNotFound` | Provider não existe | Verificar nome ou ID |
| `ErrConfigNotFound` | Chave não existe | Verificar chave ou criar config |
| `ErrGreetingNotFound` | Saudação não encontrada | Verificar keyword/language |

---

## 📄 Licença

MIT