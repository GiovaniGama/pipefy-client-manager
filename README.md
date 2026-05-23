# pipefy-client-manager

API REST de gerenciamento de clientes com integração simulada ao Pipefy via GraphQL.

## Stack

- **Go** com [Gin](https://github.com/gin-gonic/gin) (HTTP framework)
- **GORM** com SQLite (ORM + banco de dados)
- **go-playground/validator** via Gin binding (validação de inputs)
- **testify** (testes)

## Instalação e Execução

```bash
# Clone e entre no diretório
cd pipefy-client-manager

# Baixar dependências
go mod tidy

# Executar a API (porta 8080)
go run cmd/main.go
```

## Testes

```bash
go test ./tests/... -v
```

## Endpoints

### POST /clientes

Cria um novo cliente e simula a criação de um card no Pipefy.

```bash
curl -X POST http://localhost:8080/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "Atualização cadastral",
    "valor_patrimonio": 250000
  }'
```

Resposta (201):
```json
{
  "ID": 1,
  "CreatedAt": "2026-05-18T12:00:00Z",
  "nome": "João Silva",
  "email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000,
  "status": "Aguardando Análise",
  "prioridade": ""
}
```

### POST /webhooks/pipefy/card-updated

Processa o evento de atualização de card do Pipefy. Idempotente: reenviar o mesmo `event_id` retorna 200 sem reprocessar.

```bash
curl -X POST http://localhost:8080/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "evt_123",
    "card_id": "card_456",
    "cliente_email": "joao.silva@example.com",
    "timestamp": "2026-05-18T12:00:00Z"
  }'
```

Resposta (200) — primeira chamada:
```json
{
  "ID": 1,
  "status": "Processado",
  "prioridade": "prioridade_alta"
}
```

Resposta (200) — chamada duplicada:
```json
{ "message": "evento já processado" }
```

## Regra de Prioridade

| Patrimônio             | Prioridade        |
|------------------------|-------------------|
| >= R$ 200.000          | `prioridade_alta` |
| < R$ 200.000           | `prioridade_normal` |

---

## Visão de Produção (AWS)

```
                         ┌─────────────┐
Cliente / Pipefy ──────► │ API Gateway │
                         └──────┬──────┘
                                │
                    ┌───────────┴────────────┐
                    │                        │
             ┌──────▼──────┐        ┌────────▼───────┐
             │  Lambda Go  │        │   Lambda Go    │
             │ /clientes   │        │ /webhooks/...  │
             └──────┬──────┘        └────────┬───────┘
                    │                        │
          ┌─────────▼──────────┐   ┌─────────▼──────────┐
          │  RDS Aurora        │   │  RDS Aurora        │
          │  (clientes)        │   │  (clientes)        │
          └────────────────────┘   └────────┬───────────┘
                                            │
                                   ┌────────▼───────────┐
                                   │  DynamoDB          │
                                   │  (eventos -        │
                                   │   idempotência)    │
                                   └────────────────────┘
```

### Componentes

| Componente | Função |
|---|---|
| **API Gateway** | Roteamento HTTP, autenticação via API Key ou Cognito, throttling |
| **Lambda (Go binary)** | Handler stateless compilado com `GOARCH=amd64 GOOS=linux`; cada endpoint pode ser uma função separada |
| **RDS Aurora (PostgreSQL)** | Armazena a tabela `clientes`; substitui o SQLite local; suporta read replicas para escalar leituras |
| **DynamoDB** | Tabela `eventos_processados` com `event_id` como partition key e TTL configurado; acesso O(1) para verificação de idempotência sem disputar conexões do RDS |

### Por que DynamoDB para idempotência?

- Acesso em milissegundos por chave primária
- TTL nativo elimina eventos antigos automaticamente
- Sem gerenciamento de conexões (sem pool), ideal para Lambda
- Escala horizontalmente sem configuração adicional
