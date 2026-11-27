"use client";

export const SAFETY_ITEMS = [
  "DDS (Diálogo Diário de Segurança) realizado",
  "EPIs obrigatórios conferidos para toda a equipe",
  "Permissão de Trabalho (PT) emitida e válida",
  "Área de trabalho isolada e sinalizada",
  "Ferramentas e equipamentos inspecionados",
  "Riscos elétricos/altura/espaço confinado controlados",
  "Frente de trabalho limpa e organizada",
] as const;

export function SafetyChecklist({
  checked,
  onToggle,
}: {
  checked: boolean[];
  onToggle: (index: number) => void;
}) {
  return (
    <fieldset style={{ border: "none", padding: 0, margin: "12px 0" }}>
      <legend style={{ fontWeight: 600, marginBottom: 8 }}>
        Checklist de Segurança (7 itens)
      </legend>
      <div style={{ display: "grid", gap: 6 }}>
        {SAFETY_ITEMS.map((label, index) => (
          <label
            key={label}
            style={{ display: "flex", alignItems: "center", gap: 8, cursor: "pointer" }}
          >
            <input
              type="checkbox"
              checked={checked[index]}
              onChange={() => onToggle(index)}
            />
            <span>{label}</span>
          </label>
        ))}
      </div>
    </fieldset>
  );
}