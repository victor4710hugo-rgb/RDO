import type { ActivityStatus } from "@/types";

export const ACTIVITY_STATUS_LABEL: Record<ActivityStatus, string> = {
  pending: "Pendente",
  in_progress: "Em andamento",
  blocked: "Bloqueada",
  done: "Concluída",
  cancelled: "Cancelada",
};

export function State({
  loading,
  error,
  empty,
}: {
  loading: boolean;
  error: string | null;
  empty?: boolean;
}) {
  if (loading) {
    return <p className="state">Carregando…</p>;
  }

  if (error) {
    return (
      <p className="state error">
        {error}. Verifique se a API está rodando na porta 8080.
      </p>
    );
  }

  if (empty) {
    return <p className="state">Nada cadastrado ainda.</p>;
  }

  return null;
}