# RDO - Plataforma de Gestão de Obras, Serviços e Operações de Campo

Plataforma para planejamento, distribuição, execução, acompanhamento e documentação de obras, serviços e intervenções de equipes próprias e terceirizadas. O RDO (Relatório Diário de Obra) é um módulo da plataforma, não o produto inteiro.

**Status:** build beta de demonstração em desenvolvimento (Etapa 01 - Fundação). Não está pronta para produção.

## Tecnologias aprovadas

| Componente | Decisão |
|---|---|
| Frontend | Next.js + React + TypeScript |
| Backend | Go |
| Banco central | SQLite (instância única inicial) |
| Banco local | IndexedDB + Dexie.js |
| Aplicação | PWA responsiva, mobile-first |
| Offline | Requisito arquitetural completo |
| Sincronização | Fila persistente, somente com o app aberto e conexão disponível |

## Estrutura prevista

    RDO/
    |-- frontend/   (Next.js, React, TypeScript, PWA)
    |-- backend/    (API Go e SQLite)
    |-- docs/       (roteiro de testes da beta)
    |-- .gitignore
    |-- README.md

## Escopo da build beta

Entra:
- Dashboard operacional, consulta de projetos, OS e atividades.
- Formulário digital de RDO com checklist de segurança de sete itens.
- Persistência local (Dexie/IndexedDB) e fila de sincronização.
- PWA com cache dos recursos essenciais e indicadores de conexão e sincronização.

Fica para depois: autenticação real, RBAC completo, cadastros completos, planejamento automatizado, calendário, Gantt, mapas, upload de mídias para nuvem, IA e transcrição, assinaturas formais, PDFs e relatórios consolidados, hospedagem definitiva.

## Regras de integridade

- Nunca excluir dados locais porque uma tentativa de envio falhou.
- Não marcar como sincronizado antes da confirmação do servidor.
- Identificadores únicos e controle de versão; sem "última gravação vence".
- Nenhum fato é confirmado por IA sem revisão humana.

## Etapas e branches

| Etapa | Branch | Entrega |
|---|---|---|
| 01 | feat/etapa-01-fundacao | Estrutura e inicialização |
| 02 | feat/etapa-02-api-sqlite | Backend, API e banco |
| 03 | feat/etapa-03-interface | Dashboard e módulos operacionais |
| 04 | feat/etapa-04-rdo-beta | Formulário e checklist |
| 05 | feat/etapa-05-offline-sync | Dexie, fila e sincronização |
| 06 | feat/etapa-06-pwa-testes | PWA e validação integrada |

## Como executar

Será documentado nas etapas em que frontend e backend passarem a existir.