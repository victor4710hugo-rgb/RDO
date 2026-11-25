# AGENTS.md — Regras de Operação (RDO)

Leia `@ARQUITETURA.md` antes de qualquer tarefa.

## Modo de Resposta
- Entregue código pronto e direto, sem relatórios longos.
- Limite explicações a no máximo 3 tópicos curtos.
- Nunca reescreva um arquivo inteiro se apenas uma função precisa mudar.

## Arquitetura e Limites
- **Teto de 200 linhas por arquivo.**
- **Segmentação por domínio:** Ao ultrapassar 160 linhas, divida por assunto (`project_repo.go`, `os_repo.go`, `activity_repo.go`, `rdo_repo.go` e equivalentes em `models/`).
- **Fonte única nos dados:** Toda query SQL e lógica de leitura/escrita de uma entidade fica exclusivamente no seu repositório dedicado. Proibido SQL em `handlers/` ou `services/`.
- **Camadas Go:** Handler (HTTP) -> Service (Regras/Validação) -> Repository (SQL).
- **Integridade:** Nunca apague dados locais em falha de rede, nunca use "última gravação vence" e não execute `git commit`/`push` sem solicitação.