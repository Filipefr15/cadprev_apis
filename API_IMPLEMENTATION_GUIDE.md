# Guia Simplificado - Consulta de API

Este guia mostra como implementar a consulta de uma API externa de forma simples e direta.

## 📋 Passos Básicos

1. Criar a entidade (estrutura de dados)
2. Criar o cliente da API
3. Criar o handler HTTP
4. Configurar no main.go

---

## 🚀 Exemplo Prático: Consultar API de Usuários

### **Passo 1: Criar a Entidade**

📍 Arquivo: `internal/domain/entity/user.go`

```go
package entity

// User representa um usuário
type User struct {
    ID    string
    Name  string
    Email string
}
```

### **Passo 2: Criar o Cliente da API**

📍 Arquivo: `internal/api/user_client.go`

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

type UserAPIClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewUserAPIClient(baseURL string) *UserAPIClient {
    return &UserAPIClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// userResponse representa a resposta da API
type userResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// FetchUsers busca todos os usuários da API
func (c *UserAPIClient) FetchUsers(ctx context.Context) ([]*entity.User, error) {
    url := fmt.Sprintf("%s/users", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("erro ao criar requisição: %w", err)
    }
    
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("erro ao executar requisição: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status code inesperado: %d", resp.StatusCode)
    }
    
    var apiUsers []userResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiUsers); err != nil {
        return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
    }
    
    // Converter para entidades
    users := make([]*entity.User, 0, len(apiUsers))
    for _, apiUser := range apiUsers {
        users = append(users, &entity.User{
            ID:    apiUser.ID,
            Name:  apiUser.Name,
            Email: apiUser.Email,
        })
    }
    
    return users, nil
}

// FetchUserByID busca um usuário específico
func (c *UserAPIClient) FetchUserByID(ctx context.Context, id string) (*entity.User, error) {
    url := fmt.Sprintf("%s/users/%s", c.baseURL, id)
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("erro ao criar requisição: %w", err)
    }
    
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("erro ao executar requisição: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("usuário não encontrado: %s", id)
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("status code inesperado: %d", resp.StatusCode)
    }
    
    var apiUser userResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiUser); err != nil {
        return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
    }
    
    return &entity.User{
        ID:    apiUser.ID,
        Name:  apiUser.Name,
        Email: apiUser.Email,
    }, nil
}
```

### **Passo 3: Criar o Handler HTTP**

📍 Arquivo: `internal/handler/user_handler.go`

```go
package handler

import (
    "encoding/json"
    "net/http"
    
    "github.com/Filipefr15/cadprev_apis/internal/api"
)

type UserHandler struct {
    userClient *api.UserAPIClient
}

func NewUserHandler(userClient *api.UserAPIClient) *UserHandler {
    return &UserHandler{
        userClient: userClient,
    }
}

// HandleGetUsers retorna todos os usuários
func (h *UserHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.userClient.FetchUsers(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// HandleGetUser retorna um usuário específico
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
    // Extrair ID da URL (exemplo usando chi router)
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID é obrigatório", http.StatusBadRequest)
        return
    }
    
    user, err := h.userClient.FetchUserByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### **Passo 4: Configurar no main.go**

📍 Arquivo: `cmd/api/main.go`

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/Filipefr15/cadprev_apis/internal/api"
    "github.com/Filipefr15/cadprev_apis/internal/handler"
)

func main() {
    // Criar o cliente da API
    userClient := api.NewUserAPIClient("https://api.exemplo.com")
    
    // Criar o handler
    userHandler := handler.NewUserHandler(userClient)
    
    // Configurar rotas
    mux := http.NewServeMux()
    mux.HandleFunc("/users", userHandler.HandleGetUsers)
    mux.HandleFunc("/users/by-id", userHandler.HandleGetUser)
    
    // Iniciar servidor
    log.Println("Servidor rodando na porta 8080...")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}
```

---

## 📝 Resumo

É isso! Com esses 4 passos você tem:

1. ✅ Uma **entidade** para representar os dados
2. ✅ Um **cliente** que faz requisições HTTP para a API externa
3. ✅ Um **handler** que expõe endpoints HTTP
4. ✅ A **configuração** no main.go que conecta tudo

**Quando quiser adicionar persistência no banco**, você adiciona:
- Um repository (interface + implementação)
- Chama o repository.Save() depois de buscar da API

Mas por enquanto, esse código já funciona para consultar APIs externas! 🚀
