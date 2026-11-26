# AGENTS.md — Regras do Frontend (Next.js / PWA)

Siga `@../ARQUITETURA.md` e `@../AGENTS.md` antes de qualquer tarefa neste diretório.

## Modo de Resposta (Sem Enrolação)
- Entregue o código pronto para aplicar/copiar, sem relatórios longos ou teoria.
- Limite explicações a no máximo 3 tópicos curtos.
- Se precisar ler um arquivo antes de codar, leia e já entregue o código na mesma resposta.

## Regras Técnicas
- **Teto de 600 linhas por arquivo**, podendo passar sob justificativa.
- **Segmentação por domínio:** Separe tipagens e serviços por assunto quando crescerem, sempre consumindo `lib/api.ts`.
- **Rede / Android:** Proibido usar `localhost` fixo nas requisições; use sempre caminhos relativos (`/api/v1`) via `rewrites` do `next.config.ts`.
- **Offline e Sync:** Nunca exclua dados locais em falha de rede e nunca marque registro como sincronizado sem confirmação HTTP do servidor.

<!-- BEGIN:nextjs-agent-rules -->

## Atenção: Regras da Versão Instalada do Next.js

Esta versão do Next.js possui alterações incompatíveis com versões anteriores — APIs, convenções e estrutura de arquivos podem diferir do conhecimento padrão da IA. Leia o guia relevante em `node_modules/next/dist/docs/` (resolvido a partir deste diretório `frontend/`) antes de escrever qualquer código ou configurar proxy/rewrites. Respeite os avisos de depreciação.

*(Nota: Este bloco é gerenciado automaticamente pelo `next dev` via `node_modules/next/dist/server/lib/generate-agent-files.js`. Mantenha os marcadores BEGIN/END e faça commit do arquivo para evitar alterações pendentes recorrentes no Git.)*

<!-- END:nextjs-agent-rules -->