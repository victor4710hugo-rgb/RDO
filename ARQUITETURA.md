# Arquitetura — RDO

Plataforma de gestão de obras e operações de campo. Frontend Next.js (PWA) e API Go com SQLite. Build beta: offline e sincronização são requisitos, ainda incompletos.

## Árvore

```
RDO/
├── README.md
├── frontend/src/
│   ├── app/                  layout, dashboard, globals.css
│   │   ├── projetos/
│   │   ├── ordens-servico/
│   │   └── atividades/
│   ├── components/           State.tsx, layout/ (AppShell, Sidebar, Topbar)
│   ├── lib/api.ts
│   └── types/index.ts
├── backend/
│   ├── cmd/api/main.go
│   └── internal/
│       ├── router/           rotas HTTP
│       ├── handlers/         health, operations, field_reports, sync
│       ├── services/         regra do RDO
│       ├── repositories/     SQL
│       ├── models/           structs de domínio
│       └── database/         sqlite, migrations, seed
└── Docs/                     propostas visuais
```
## Arquivos principais

- `backend/cmd/api/main.go` — abre o SQLite, migra, semeia e sobe a API em `:8080`.
- `internal/router/router.go` — registra `/api/v1` e injeta repositório e serviços.
- `handlers/*.go` — HTTP: leem o pedido, chamam a camada de baixo e devolvem JSON.
- `services/field_report_service.go` — valida o RDO (id, OS, data, payload) antes de gravar.
- `repositories/repository.go` — único acesso SQL (projetos, OS, atividades, RDO).
- `models/models.go` — structs de domínio e JSON.
- `database/sqlite.go`, `migrations.go`, `seed.go` — conexão, schema e dados iniciais.
- `frontend/src/app/layout.tsx` — shell (sidebar e topbar) em torno das páginas.
- `page.tsx` e as pastas de projetos, OS e atividades — dashboard e listas.
- `lib/api.ts` — cliente HTTP e hook `useApi`.
- `types/index.ts` — tipos compartilhados do frontend.
- `components/State.tsx` — carregamento, erro e vazio.
## Regras

- Nunca misturar acesso ao banco dentro das rotas HTTP. SQL fica em `repositories`; regra de negócio em `services`; o handler só traduz HTTP.
- Arquivos com no máximo 200 linhas. Se o arquivo passar disso justificar.
- Nunca excluir dados locais porque um envio falhou.
- Não marcar como sincronizado antes da confirmação do servidor.
- Identificadores únicos e controle de versão; sem "última gravação vence".
- Nenhum fato é confirmado por IA sem revisão humana.
