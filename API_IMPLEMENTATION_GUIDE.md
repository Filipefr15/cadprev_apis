# Guia de Implementação de Consumo de APIs
## Padrões DDD + Clean Architecture

Este guia explica como implementar o consumo de APIs externas mantendo os padrões de **Domain-Driven Design (DDD)** e **Clean Architecture** já estabelecidos neste repositório.

---

## 📋 Índice

1. [Visão Geral da Arquitetura](#visão-geral-da-arquitetura)
2. [Estrutura de Camadas](#estrutura-de-camadas)
3. [Passo a Passo para Nova API](#passo-a-passo-para-nova-api)
4. [Exemplos Práticos](#exemplos-práticos)
5. [Boas Práticas](#boas-práticas)
6. [Checklist de Implementação](#checklist-de-implementação)

---

## 🏗️ Visão Geral da Arquitetura

Este projeto segue o princípio de **Dependency Inversion** (SOLID), onde:

- **Camadas internas** (Domain) definem **interfaces**
- **Camadas externas** (Infrastructure) **implementam** essas interfaces
- As dependências apontam **sempre para dentro** (para o Domain)

```
┌─────────────────────────────────────────┐
│         cmd/api (Entry Point)           │
│         - main.go                       │
│         - Dependency Injection          │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Handler Layer                   │
│         (Presentation/Interface)        │
│         - HTTP Handlers                 │
│         - Request/Response DTOs         │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Use Case Layer                  │
│         (Application/Business Logic)    │
│         - Orchestration                 │
│         - Business Rules                │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Domain Layer (Core)             │
│         - Entities                      │
│         - Interfaces (Ports)            │
│         - Business Logic                │
└─────────────────────────────────────────┘
               ▲
               │
┌──────────────┴──────────────────────────┐
│         Infrastructure Layer            │
│         - API Clients                   │
│         - Repositories                  │
│         - External Integrations         │
└─────────────────────────────────────────┘
```

---

## 📂 Estrutura de Camadas

### 1. **Domain Layer** (`internal/domain/`)

**Propósito**: Núcleo da aplicação, sem dependências externas

```
internal/domain/
├── entity/           # Entidades de domínio
│   └── data.go      # Entidade Data
├── api_client.go    # Interface para clientes de API
└── repository.go    # Interface para repositórios
```

**Características:**
- ✅ Define **interfaces** (contratos)
- ✅ Contém **entidades** com lógica de negócio
- ❌ **NÃO** importa pacotes de infraestrutura
- ❌ **NÃO** sabe como os dados são buscados

### 2. **Infrastructure Layer** (`internal/api/` e `internal/repository/`)

**Propósito**: Implementações concretas das interfaces do domain

```
internal/
├── api/
│   └── stub_client.go    # Implementação do cliente API
└── repository/
    └── inmemory.go        # Implementação do repositório
```

**Características:**
- ✅ **Implementa** interfaces do domain
- ✅ Lida com detalhes técnicos (HTTP, JSON, DB)
- ✅ Pode usar bibliotecas externas

### 3. **Use Case Layer** (`internal/usecase/`)

**Propósito**: Orquestração de lógica de negócio

```
internal/usecase/
└── ingest.go    # Caso de uso para ingestão de dados
```

**Características:**
- ✅ Orquestra operações entre domain e infrastructure
- ✅ Implementa regras de negócio de aplicação
- ✅ Depende apenas de **interfaces** do domain

### 4. **Handler Layer** (`internal/handler/`)

**Propósito**: Interface HTTP/REST da aplicação

```
internal/handler/
└── data_handler.go    # Handler HTTP
```

**Características:**
- ✅ Recebe requisições HTTP
- ✅ Valida entrada
- ✅ Chama use cases
- ✅ Retorna respostas HTTP

### 5. **Entry Point** (`cmd/api/`)

**Propósito**: Inicialização e injeção de dependências

```
cmd/api/
└── main.go    # Ponto de entrada
```

**Características:**
- ✅ **Dependency Injection** manual
- ✅ Configuração da aplicação
- ✅ Inicialização do servidor

---

## 🚀 Passo a Passo para Nova API

### **Exemplo**: Implementar consumo da API de Usuários

#### **Passo 1: Definir a Entidade no Domain**

📍 Local: `internal/domain/entity/user.go`

```go
package entity

import "time"

// User representa a entidade de usuário do domínio
type User struct {
    ID        string
    Name      string
    Email     string
    Role      string
    Active    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewUser cria um novo usuário com validações
func NewUser(id, name, email, role string) (*User, error) {
    // Validações de negócio aqui
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }
    if email == "" {
        return nil, fmt.Errorf("email cannot be empty")
    }
    
    now := time.Now()
    return &User{
        ID:        id,
        Name:      name,
        Email:     email,
        Role:      role,
        Active:    true,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
    u.Active = false
    u.UpdatedAt = time.Now()
}
```

#### **Passo 2: Definir Interface no Domain**

📍 Local: `internal/domain/user_api_client.go`

```go
package domain

import (
    "context"
    
    "github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// UserAPIClient define a interface para comunicação com API de usuários
// Esta interface pertence à camada de domínio mas é implementada na infraestrutura
type UserAPIClient interface {
    FetchUsers(ctx context.Context) ([]*entity.User, error)
    FetchUserByID(ctx context.Context, id string) (*entity.User, error)
    CreateUser(ctx context.Context, user *entity.User) error
}
```

#### **Passo 3: Definir Interface de Repositório (se necessário)**

📍 Local: `internal/domain/user_repository.go`

```go
package domain

import (
    "context"
    
    "github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// UserRepository define a interface para persistência de usuários
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context) ([]*entity.User, error)
    Delete(ctx context.Context, id string) error
}
```

#### **Passo 4: Implementar Cliente da API (Infrastructure)**

📍 Local: `internal/api/user_client.go`

```go
package api

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// UserAPIClient implementa a interface domain.UserAPIClient
type UserAPIClient struct {
    baseURL    string
    httpClient *http.Client
}

// NewUserAPIClient cria um novo cliente da API de usuários
func NewUserAPIClient(baseURL string) *UserAPIClient {
    return &UserAPIClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// userResponse representa a resposta da API externa
// Este DTO pertence à camada de infraestrutura
type userResponse struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    Active bool   `json:"active"`
}

// FetchUsers busca todos os usuários da API externa
func (c *UserAPIClient) FetchUsers(ctx context.Context) ([]*entity.User, error) {
    url := fmt.Sprintf("%s/users", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Adicionar headers se necessário
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    var apiUsers []userResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiUsers); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    // Converter DTOs da API para entidades de domínio
    users := make([]*entity.User, 0, len(apiUsers))
    for _, apiUser := range apiUsers {
        user, err := entity.NewUser(apiUser.ID, apiUser.Name, apiUser.Email, apiUser.Role)
        if err != nil {
            // Log e continue ou retorne erro dependendo do requisito
            continue
        }
        users = append(users, user)
    }
    
    return users, nil
}

// FetchUserByID busca um usuário específico por ID
func (c *UserAPIClient) FetchUserByID(ctx context.Context, id string) (*entity.User, error) {
    url := fmt.Sprintf("%s/users/%s", c.baseURL, id)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("user not found: %s", id)
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }
    
    var apiUser userResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiUser); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return entity.NewUser(apiUser.ID, apiUser.Name, apiUser.Email, apiUser.Role)
}

// CreateUser cria um novo usuário na API externa
func (c *UserAPIClient) CreateUser(ctx context.Context, user *entity.User) error {
    // Implementação similar...
    return nil
}
```

#### **Passo 5: Implementar Use Case**

📍 Local: `internal/usecase/user_sync.go`

```go
package usecase

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Filipefr15/cadprev_apis/internal/domain"
)

// UserSyncUseCase orquestra a sincronização de usuários
type UserSyncUseCase struct {
    userAPIClient domain.UserAPIClient
    userRepo      domain.UserRepository
}

// NewUserSyncUseCase cria um novo use case com injeção de dependências
func NewUserSyncUseCase(
    userAPIClient domain.UserAPIClient,
    userRepo domain.UserRepository,
) *UserSyncUseCase {
    return &UserSyncUseCase{
        userAPIClient: userAPIClient,
        userRepo:      userRepo,
    }
}

// SyncAllUsers busca usuários da API e persiste no repositório
func (uc *UserSyncUseCase) SyncAllUsers(ctx context.Context) error {
    // Buscar usuários da API externa
    users, err := uc.userAPIClient.FetchUsers(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch users from API: %w", err)
    }
    
    log.Printf("Fetched %d users from external API", len(users))
    
    // Persistir cada usuário
    successCount := 0
    for _, user := range users {
        if err := uc.userRepo.Save(ctx, user); err != nil {
            log.Printf("Failed to save user %s: %v", user.ID, err)
            continue
        }
        successCount++
    }
    
    log.Printf("Successfully saved %d/%d users", successCount, len(users))
    return nil
}

// SyncUserByID sincroniza um usuário específico
func (uc *UserSyncUseCase) SyncUserByID(ctx context.Context, id string) error {
    user, err := uc.userAPIClient.FetchUserByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to fetch user: %w", err)
    }
    
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    log.Printf("Successfully synced user: %s", user.ID)
    return nil
}
```

#### **Passo 6: Implementar Handler HTTP**

📍 Local: `internal/handler/user_handler.go`

```go
package handler

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    
    "github.com/Filipefr15/cadprev_apis/internal/usecase"
)

// UserHandler lida com requisições HTTP relacionadas a usuários
type UserHandler struct {
    syncUseCase *usecase.UserSyncUseCase
}

// NewUserHandler cria um novo handler com injeção de dependências
func NewUserHandler(syncUseCase *usecase.UserSyncUseCase) *UserHandler {
    return &UserHandler{
        syncUseCase: syncUseCase,
    }
}

// SyncAllUsers lida com POST /api/users/sync
func (h *UserHandler) SyncAllUsers(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    ctx := r.Context()
    
    err := h.syncUseCase.SyncAllUsers(ctx)
    if err != nil {
        log.Printf("Error syncing users: %v", err)
        http.Error(w, "Failed to sync users", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Users synchronized successfully",
    })
}

// SyncUserByID lida com POST /api/users/sync/{id}
func (h *UserHandler) SyncUserByID(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Extrair ID da URL
    path := strings.TrimPrefix(r.URL.Path, "/api/users/sync/")
    id := strings.TrimSuffix(path, "/")
    
    if id == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }
    
    ctx := r.Context()
    
    err := h.syncUseCase.SyncUserByID(ctx, id)
    if err != nil {
        log.Printf("Error syncing user %s: %v", id, err)
        http.Error(w, "Failed to sync user", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "User synchronized successfully",
        "id":      id,
    })
}

// RegisterRoutes registra as rotas do handler
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/users/sync", h.SyncAllUsers)
    mux.HandleFunc("/api/users/sync/", h.SyncUserByID)
}
```

#### **Passo 7: Configurar Dependency Injection**

📍 Local: `cmd/api/main.go`

```go
func main() {
    // ... código existente ...
    
    // Nova API - User Client
    userAPIURL := getEnv("USER_API_URL", "https://api.example.com")
    userAPIClient := api.NewUserAPIClient(userAPIURL)
    log.Printf("Initialized user API client with base URL: %s", userAPIURL)
    
    // User Repository
    userRepository := repository.NewInMemoryUserRepository()
    log.Println("Initialized user repository")
    
    // User Use Case
    userSyncUseCase := usecase.NewUserSyncUseCase(userAPIClient, userRepository)
    log.Println("Initialized user sync use case")
    
    // User Handler
    userHandler := handler.NewUserHandler(userSyncUseCase)
    log.Println("Initialized user handler")
    
    // Registrar rotas
    mux := http.NewServeMux()
    dataHandler.RegisterRoutes(mux)
    userHandler.RegisterRoutes(mux)  // Nova rota
    
    // ... resto do código ...
}
```

---

## 📚 Exemplos Práticos

### **Exemplo 1: API REST Simples (GET)**

```go
// internal/api/product_client.go
func (c *ProductAPIClient) FetchProducts(ctx context.Context) ([]*entity.Product, error) {
    url := fmt.Sprintf("%s/products", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()
    
    // ... processamento ...
}
```

### **Exemplo 2: API com Autenticação**

```go
// internal/api/auth_client.go
type AuthenticatedClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func (c *AuthenticatedClient) FetchData(ctx context.Context) ([]*entity.Data, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
    if err != nil {
        return nil, err
    }
    
    // Adicionar autenticação
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
    req.Header.Set("X-API-Key", c.apiKey)
    
    resp, err := c.httpClient.Do(req)
    // ... resto da implementação ...
}
```

### **Exemplo 3: API com Paginação**

```go
func (c *PaginatedAPIClient) FetchAllData(ctx context.Context) ([]*entity.Data, error) {
    var allData []*entity.Data
    page := 1
    pageSize := 100
    
    for {
        url := fmt.Sprintf("%s/data?page=%d&size=%d", c.baseURL, page, pageSize)
        
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return nil, err
        }
        
        resp, err := c.httpClient.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        
        var pageData []entity.Data
        if err := json.NewDecoder(resp.Body).Decode(&pageData); err != nil {
            return nil, err
        }
        
        if len(pageData) == 0 {
            break // Última página
        }
        
        allData = append(allData, pageData...)
        page++
    }
    
    return allData, nil
}
```

### **Exemplo 4: API com Retry e Circuit Breaker**

```go
import "github.com/sony/gobreaker"

type ResilientAPIClient struct {
    baseURL       string
    httpClient    *http.Client
    circuitBreaker *gobreaker.CircuitBreaker
}

func NewResilientAPIClient(baseURL string) *ResilientAPIClient {
    settings := gobreaker.Settings{
        Name:        "API Client",
        MaxRequests: 3,
        Timeout:     60 * time.Second,
    }
    
    return &ResilientAPIClient{
        baseURL:       baseURL,
        httpClient:    &http.Client{Timeout: 30 * time.Second},
        circuitBreaker: gobreaker.NewCircuitBreaker(settings),
    }
}

func (c *ResilientAPIClient) FetchData(ctx context.Context) ([]*entity.Data, error) {
    result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
        return c.doFetchData(ctx)
    })
    
    if err != nil {
        return nil, err
    }
    
    return result.([]*entity.Data), nil
}
```

---

## ✅ Boas Práticas

### **1. Separação de Responsabilidades**

❌ **Errado**: Domain depende de infraestrutura
```go
// internal/domain/user.go
import "net/http"  // ❌ Domain não deve importar pacotes de infraestrutura

type User struct {
    httpClient *http.Client  // ❌
}
```

✅ **Correto**: Domain define interface, infraestrutura implementa
```go
// internal/domain/user_api_client.go
type UserAPIClient interface {
    FetchUsers(ctx context.Context) ([]*entity.User, error)
}

// internal/api/user_client.go
import "net/http"  // ✅ Infraestrutura pode usar HTTP

type UserAPIClient struct {
    httpClient *http.Client  // ✅
}
```

### **2. DTOs vs Entidades**

❌ **Errado**: Usar estrutura da API diretamente no domain
```go
// ❌ Entidade acoplada ao formato da API
type User struct {
    UserID    string `json:"user_id"`     // Campo específico da API
    FullName  string `json:"full_name"`   // Nome específico da API
}
```

✅ **Correto**: DTOs na infraestrutura, entidades no domain
```go
// internal/domain/entity/user.go - Entidade pura
type User struct {
    ID   string
    Name string
}

// internal/api/user_client.go - DTO específico da API
type userAPIResponse struct {
    UserID   string `json:"user_id"`
    FullName string `json:"full_name"`
}

// Conversão de DTO para Entidade
func (r *userAPIResponse) ToEntity() *entity.User {
    return &entity.User{
        ID:   r.UserID,
        Name: r.FullName,
    }
}
```

### **3. Context para Cancelamento**

✅ **Sempre** passe `context.Context` para operações de rede
```go
func (c *APIClient) FetchData(ctx context.Context) ([]*entity.Data, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
    // ...
}
```

### **4. Tratamento de Erros**

✅ **Use** `fmt.Errorf` com `%w` para wrapping de erros
```go
if err != nil {
    return nil, fmt.Errorf("failed to fetch users: %w", err)
}
```

### **5. Logging Apropriado**

```go
// Use case layer
log.Printf("Fetched %d records from external API", len(dataList))

// Handler layer
log.Printf("Error syncing users: %v", err)
```

### **6. Timeouts e Retries**

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

### **7. Testes Unitários**

```go
// internal/api/user_client_test.go
func TestUserAPIClient_FetchUsers(t *testing.T) {
    // Mock do servidor HTTP
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode([]userResponse{
            {ID: "1", Name: "John", Email: "john@example.com"},
        })
    }))
    defer server.Close()
    
    client := NewUserAPIClient(server.URL)
    users, err := client.FetchUsers(context.Background())
    
    assert.NoError(t, err)
    assert.Len(t, users, 1)
    assert.Equal(t, "John", users[0].Name)
}
```

---

## ✅ Checklist de Implementação

Use esta checklist ao implementar uma nova API:

### **Domain Layer**
- [ ] Criar entidade em `internal/domain/entity/`
- [ ] Adicionar validações de negócio na entidade
- [ ] Definir interface do cliente em `internal/domain/`
- [ ] Definir interface do repositório (se necessário)

### **Infrastructure Layer**
- [ ] Criar cliente da API em `internal/api/`
- [ ] Definir DTOs específicos da API (não expor ao domain)
- [ ] Implementar conversão DTO → Entidade
- [ ] Adicionar tratamento de erros HTTP
- [ ] Configurar timeout e retry (se necessário)
- [ ] Implementar autenticação (se necessário)

### **Use Case Layer**
- [ ] Criar use case em `internal/usecase/`
- [ ] Injetar dependências via interfaces do domain
- [ ] Implementar lógica de orquestração
- [ ] Adicionar logs apropriados

### **Handler Layer**
- [ ] Criar handler em `internal/handler/`
- [ ] Implementar endpoints HTTP
- [ ] Validar request/response
- [ ] Adicionar método `RegisterRoutes`

### **Entry Point**
- [ ] Adicionar variável de ambiente para URL da API
- [ ] Instanciar cliente da API
- [ ] Instanciar repositório (se necessário)
- [ ] Instanciar use case
- [ ] Instanciar handler
- [ ] Registrar rotas

### **Testes**
- [ ] Testes unitários do cliente da API
- [ ] Testes unitários do use case
- [ ] Testes de integração (opcional)

### **Documentação**
- [ ] Atualizar README com novos endpoints
- [ ] Documentar variáveis de ambiente
- [ ] Adicionar exemplos de uso

---

## 🎯 Princípios SOLID Aplicados

### **S - Single Responsibility**
- Cada camada tem uma responsabilidade clara
- Cliente de API: buscar dados
- Use Case: orquestração
- Handler: interface HTTP

### **O - Open/Closed**
- Aberto para extensão (novas APIs)
- Fechado para modificação (interfaces estáveis)

### **L - Liskov Substitution**
- Qualquer implementação de `ExternalAPIClient` pode ser substituída
- Stub em desenvolvimento, cliente real em produção

### **I - Interface Segregation**
- Interfaces pequenas e específicas
- `UserAPIClient` separado de `DataAPIClient`

### **D - Dependency Inversion**
- Use cases dependem de **interfaces** (domain)
- Não dependem de **implementações** (infrastructure)

---

## 🔄 Fluxo de Dados

```
┌─────────────┐
│ HTTP Request│
└──────┬──────┘
       │
       ▼
┌──────────────┐
│   Handler    │ ← Valida request
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Use Case    │ ← Orquestra lógica de negócio
└──────┬───────┘
       │
       ▼
┌──────────────────────────┐
│  API Client (Interface)  │ ← Interface do domain
└──────┬───────────────────┘
       │
       ▼
┌──────────────────────────┐
│ API Client (Implementação)│ ← Faz requisição HTTP
└──────┬───────────────────┘
       │
       ▼
┌──────────────┐
│ External API │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ DTO Response │ ← Converte para entidade
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Entity     │ ← Retorna para use case
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │ ← Persiste (se necessário)
└──────────────┘
```

---

## 📖 Referências

- **Clean Architecture**: Robert C. Martin
- **Domain-Driven Design**: Eric Evans
- **SOLID Principles**: Robert C. Martin
- **Go Standard Library**: `net/http`, `context`
- **Dependency Injection Pattern**: Martin Fowler

---

## 🚨 Antipadrões a Evitar

### ❌ **1. God Object**
```go
// ❌ Classe faz tudo
type Service struct {
    // Busca, processa, salva, envia email, etc.
}
```

### ❌ **2. Acoplamento Direto**
```go
// ❌ Use case depende de implementação concreta
func NewUserUseCase(client *api.UserAPIClient) *UserUseCase {
    // Deveria depender de interface
}
```

### ❌ **3. Lógica de Negócio no Handler**
```go
// ❌ Handler com lógica de negócio
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // Validação complexa
    // Cálculos de negócio
    // Transformações
    // Isso deve estar no use case!
}
```

### ❌ **4. Entidades Anêmicas**
```go
// ❌ Entidade sem comportamento
type User struct {
    ID   string
    Name string
    // Apenas getters/setters, sem lógica de negócio
}
```

---

## 🎓 Conclusão

Este guia estabelece os padrões para implementação de consumo de APIs mantendo:

✅ **Separação clara de responsabilidades**  
✅ **Testabilidade** (interfaces facilitam mocks)  
✅ **Manutenibilidade** (mudanças isoladas por camada)  
✅ **Escalabilidade** (fácil adicionar novas APIs)  
✅ **Princípios SOLID e DDD**

Para dúvidas ou sugestões, consulte a equipe de arquitetura do projeto.
