# RDO — Gestão de Obras e Operações de Campo

Aplicação web em desenvolvimento para consultar projetos, ordens de serviço (OS) e atividades de operações de campo. O Relatório Diário de Obra (RDO) faz parte do domínio e do schema do banco, mas o fluxo completo de preenchimento e sincronização ainda não está implementado.

> **Estado atual:** protótipo com dashboard e telas de consulta no frontend e uma API Go em desenvolvimento. A API não compila no estado atual: o código importa `golang.org/x/crypto/bcrypt`, mas esse módulo não está declarado em `backend/go.mod`. As rotas de sincronização ainda são respostas de status, e não sincronizam dados.

## Funcionalidades implementadas

- Dashboard com totais de projetos, OS e atividades, além do percentual de atividades concluídas.
- Consulta de projetos, com pesquisa local pelo nome do projeto ou cliente.
- Consulta de ordens de serviço com progresso calculado a partir das atividades concluídas.
- Consulta de atividades e indicação de status.
- API HTTP JSON para consultar projetos, OS e atividades.
- Banco SQLite com schema versionado, chaves estrangeiras e modo WAL; na primeira execução, são inseridos dados operacionais de demonstração se ainda não houver projetos.

O schema prevê outros domínios (usuários, equipes, RDO, checklists e sincronização), mas isso não significa que todos tenham telas ou operações disponíveis. O serviço de usuários não está conectado às rotas da API. Não há, no frontend atual, formulário de RDO, autenticação, persistência offline com IndexedDB, PWA ou sincronização funcional. As rotas de `sync` retornam somente `{"status":"ok"}`.

## Tecnologias

- **Frontend:** Next.js 16 (App Router), React 19 e TypeScript 5.
- **Backend:** Go (versão declarada no `go.mod`: 1.27) e `net/http`.
- **Banco:** SQLite, usando `modernc.org/sqlite`; UUIDs com `github.com/google/uuid`.
- **Qualidade do frontend:** ESLint 9 e configuração `eslint-config-next`.

O código do backend também usa `golang.org/x/crypto/bcrypt`, mas essa dependência está ausente do `go.mod` e atualmente impede o teste/build do backend.

## Estrutura do projeto

```text
.
├── README.md
├── ARQUITETURA.md
├── rdo_etapa_d2_modelo_e_especificacao.md
├── Docs/                         # Propostas visuais
├── backend/
│   ├── cmd/api/main.go           # Inicialização da API
│   ├── data/                     # Banco SQLite criado em execução
│   └── internal/
│       ├── database/             # Conexão, schema, migrations e seed
│       ├── handlers/             # Handlers HTTP
│       ├── models/               # Estruturas e JSON
│       ├── repositories/         # Consultas e gravações SQLite
│       ├── router/               # Registro das rotas
│       └── services/             # Regras de negócio
└── frontend/
    ├── package.json              # Scripts e dependências npm
    └── src/
        ├── app/                  # Dashboard e telas de consulta
        ├── components/           # Shell, navegação e estados de interface
        ├── lib/api.ts            # Cliente HTTP do frontend
        └── types/                # Tipos TypeScript
```

## Pré-requisitos e instalação

- Go compatível com a versão declarada em `backend/go.mod`.
- Node.js e npm compatíveis com as dependências do frontend.

Instale as dependências do frontend:

```sh
cd frontend
npm ci
```

Baixe as dependências Go:

```sh
cd backend
go mod download
```

Não há arquivo `.env.example`, `Dockerfile` ou configuração de deploy no repositório. Não versione credenciais.

## Variáveis de ambiente

| Variável | Uso no código |
|---|---|
| `NEXT_PUBLIC_API_URL` | URL base da API no frontend. Se não for definida, usa `http://localhost:8080/api/v1`. |
| `ADMIN_EMAIL` | Lida pelo seed do banco junto com `ADMIN_PASSWORD`. O seed tenta criar um administrador, mas o papel `admin` não é inserido pelo seed atual, então essa configuração pode fazer a inicialização falhar. Também é consultada pelo serviço de usuários, que não está ligado à inicialização atual. |
| `ADMIN_PASSWORD` | Senha que aciona o seed de administrador junto com `ADMIN_EMAIL`; consulte a limitação descrita acima. |
| `ADMIN_INITIAL_PASSWORD` | Consultada pelo serviço de usuários para criação inicial de usuário; esse serviço não é chamado pelo fluxo atual de inicialização da API. |
| `MANAGER_EMAIL` | Consultada pelo serviço de usuários, atualmente não conectado à API. |
| `MANAGER_INITIAL_PASSWORD` | Consultada pelo serviço de usuários, atualmente não conectado à API. |

O frontend e o backend não carregam automaticamente um arquivo `.env` por configuração própria. Defina as variáveis no ambiente do processo que iniciar cada aplicação.

## Como executar

### Frontend — desenvolvimento

```sh
cd frontend
npm run dev
```

Abra <http://localhost:3000>. O cliente do frontend consulta a API em `http://localhost:8080/api/v1` por padrão. A API não configura CORS; portanto, o acesso entre essas origens pode ser bloqueado pelo navegador até que o backend ou a configuração local trate essa política.

### Backend — desenvolvimento

```sh
cd backend
go run ./cmd/api
```

Quando iniciar, a API escuta em `:8080` e cria/atualiza `backend/data/rdo.db` com migrations e dados de demonstração. **Esse comando está atualmente bloqueado** pela dependência `golang.org/x/crypto/bcrypt` ausente em `go.mod`; o teste abaixo reproduz o problema.

### Build e verificações

```sh
# Frontend
cd frontend
npm run lint
npx tsc --noEmit
npm run build
npm run start       # após o build

# Backend (a partir de backend/)
go test ./...
go build ./cmd/api
```

Não há testes automatizados implementados nos pacotes Go. No estado verificado, `npm run lint` termina sem erros, com um aviso de dependência do hook `useEffect` em `frontend/src/lib/api.ts`; `go test ./...` falha porque `golang.org/x/crypto/bcrypt` não está declarado no módulo.

## API e uso

As rotas registradas usam o prefixo `/api/v1`:

| Método | Rota | Comportamento atual |
|---|---|---|
| `GET` | `/health` | Retorna `{"status":"ok"}`. |
| `GET` | `/projects` | Lista projetos ativos com o nome do cliente. |
| `GET` | `/work-orders` | Lista ordens de serviço ativas. |
| `GET` | `/activities` | Lista atividades ativas. |
| `POST` | `/field-reports` | Rota de gravação registrada, mas ainda incompleta e não validada como fluxo funcional. |
| `GET` | `/sync/pull` | Retorna `{"status":"ok"}`; não busca nem transfere registros. |
| `POST` | `/sync/push` | Retorna `{"status":"ok"}`; não processa nem persiste registros. |

Exemplos de consulta, após a API estar executável:

```sh
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/projects
curl http://localhost:8080/api/v1/work-orders
curl http://localhost:8080/api/v1/activities
```

## Escopo ainda indefinido

- Autenticação e autorização disponíveis pela API: **não implementadas**.
- Formulário de RDO, checklist no frontend e sincronização offline: **não implementados**.
- Deploy e licença: **a definir**.
