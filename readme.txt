# gamestats-backend

Backend em Go projetado para ingestão, indexação e consulta de estatísticas de partidas em servidores Arma Reforger-like. Ele coleta eventos (principalmente kills), gera métricas por jogador, rounds, servidores, e oferece endpoints REST otimizados para dashboards, bots e integrações.

Este README descreve completamente a arquitetura, fluxo interno, pontos fracos atuais e melhorias planejadas - incluindo a implementação futura de **caching com Redis**.

---

## 🚀 Objetivos do Projeto

* Processar eventos de jogo em **alta velocidade**.
* Centralizar estatísticas de players, rounds, armas e servidores.
* Servir dados agregados de forma consistente e de baixa latência.
* Facilitar futuras integrações com dashboards web, bots Discord, API pública e pipelines ETL.

---

## 📂 Estrutura do Projeto

```
.
├── cmd/
│   └── main.go                # bootstrap da aplicação
├── internal/
│   ├── controllers/           # handlers HTTP (Gin)
│   ├── services/              # regra de negócio
│   ├── repositories/          # GORM + SQL
│   ├── entities/              # modelos do banco
│   ├── dtos/                  # payloads e validações
│   └── database/              # init do DB
├── migrations/                # scripts SQL (PostgreSQL)
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

---

## ⚙️ Requisitos

* **Go 1.24+**
* **PostgreSQL 15+**
* **Docker e Docker Compose** (recomendado)
* (futuro) **Redis** para cache

---

## 🔧 Configuração - Variáveis de ambiente

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Os principais valores:

```
POSTGRESQL_USERNAME=
POSTGRESQL_PASSWORD=
POSTGRESQL_GAMESTATS_DATABASE=
POSTGRESQL_HOST=
POSTGRESQL_PORT=5432
POSTGRESQL_SSLMODE=disable

API_KEY=
APP_ENV=development
GIN_PORT=8080
```

Se `API_KEY` estiver setado, o middleware de autenticação é ativado.

---

## 🗄️ Migrações

Os arquivos SQL estão em `/migrations`. Como são scripts PostgreSQL, é recomendável executar com Flyway ou psql.

Exemplo:

```bash
psql $DATABASE_URL -f migrations/V1__create_init_schema.sql
```

As tabelas principais criadas são:

* `kills`
* `players`
* `servers`
* `rounds`
* `players_stats`

---

## ▶️ Executando Localmente

```bash
go mod download
go run ./cmd
```

API disponível em:

```
http://localhost:8080
```

---

## 🐳 Executando com Docker

```bash
docker compose up --build
```

O docker-compose inicia API + PostgreSQL automaticamente.

---

## 📡 Endpoints Principais

Base: `/api/v1`

### 🔫 Events - Kills

* `POST /api/v1/events/kill`

Aceita **array** de kills:

```json
[
  {
    "server_id": "uuid",
    "round_id": "uuid",
    "killer_id": "uuid",
    "victim_id": "uuid",
    "distance": 23.4,
    "is_headshot": true,
    "is_friendly": false,
    "timestamp": "2025-12-11T10:00:00Z"
  }
]
```

### 🧍 Players

* `GET /api/v1/players/:id`
* `GET /api/v1/players/:id/stats`
* `GET /api/v1/players/leaderboard`

### 🏳️ Servers

* CRUD completo
* `GET /api/v1/servers/:id/rounds`

### 🎯 Rounds

* criação, finalização, scoreboard, histórico por player.

---

## 🧠 Arquitetura Interna

Fluxo simplificado:

```
JSON -> DTO -> Controller -> Service -> Repository -> PostgreSQL
                          ^
                          | -> Redis Cache
```

A aplicação segue camadas claras:

* **DTOs** validam payload.
* **Controllers** convertem erros, padrões de resposta e status HTTP.
* **Services** aplicam lógica de agregação e cálculo estatístico.
* **Repositories** executam queries otimizadas.

---

## ⚠️ Pontos que Necessitam Melhoria

### 1. **Falta de caching -> alta carga em consultas populares**

Leaderboards, estatísticas agregadas e scoreboards dependem de queries pesadas (JOINs + GROUP BY). Isso escalará mal sem cache.

### 2. **Ausência de transações em algumas operações**

Salvar kills em batch deveria utilizar `BEGIN/COMMIT` para garantir atomicidade.

### 3. **Camada de serviços ainda mistura lógica com detalhes do repositório**

Seria ideal isolar completamente cálculos de hits/kills/historic.

### 4. **GORM pode gerar queries subótimas**

Algumas rotas se beneficiariam de SQL manual.

### 5. **Falta de rate limit e pagination robusta**

Rotas como leaderboard podem retornar payloads muito grandes.

### 6. **Ausência de testes para controllers com Gin**

Atualmente os testes focam mais nos repositórios.

---

## 🔮 Melhorias Futuras (Roadmap)

### ✔️ **Implementar Redis para caching inteligente**

Sugestão para o cache:

| Endpoint               | Tipo            | TTL    | Observação                          |
| ---------------------- | --------------- | ------ | ----------------------------------- |
| Leaderboard            | FULL PAGE CACHE | 30s    | alta demanda                        |
| Player Stats           | OBJECT CACHE    | 10–20s | muda pouco                          |
| Round Scoreboard       | OBJECT CACHE    | 60s    | somente enquanto o round está ativo |
| Servers List           | STATIC CACHE    | 5m     | baixo churn                         |
| Kills Streaming        | OBJECT CACHE    | 10s    | atualizações em tempo real          |

#### Estratégia recomendada

1. Redis como write-through cache.
2. Chaves com prefixos:

   * `leaderboard:global`
   * `player:stats:<id>`
   * `round:scoreboard:<roundId>`
3. Invalidação automática quando existem novos eventos.
4. Usar Redis Cluster para escalabilidade futura.

### ✔️ Adicionar suporte a WebSocket

Para atualizar dashboards em tempo real.

### ✔️ Criar job async para agregações

Atualmente algumas métricas são calculadas on-demand.

### ✔️ Logging estruturado + tracing (OpenTelemetry)

Simplifica debugging e análise de performance.

### ✔️ Suporte a sharding de banco em múltiplos servidores

Para lidar com grandes volumes de kills.

---

## 🧪 Testes

Rodar testes:

```bash
go test ./... -v
```

Recomendado adicionar:

* mocks de Redis
* testes de integração para batches de kills
* testes de concorrência

---
