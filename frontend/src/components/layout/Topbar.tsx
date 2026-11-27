"use client";

import { usePathname } from "next/navigation";

const ROUTE_TITLES: Record<string, string> = {
  "/": "Dashboard",
  "/projetos": "Projetos",
  "/ordens-servico": "Ordens de Serviço",
  "/atividades": "Atividades",
  "/rdo": "Relatórios Diários (RDO)",
};

export default function Topbar({ onMenuToggle }: { onMenuToggle: () => void }) {
  const pathname = usePathname();
  const current = ROUTE_TITLES[pathname] ?? "Operações";

  function toggleTheme() {
    const root = document.documentElement;
    root.dataset.theme = root.dataset.theme === "dark" ? "light" : "dark";
  }

  return (
    <header className="topbar">
      <button type="button" className="mb" onClick={onMenuToggle}>
        Menu
      </button>
      <span className="crumb">
        Plataforma RDO › <b>{current}</b>
      </span>
      <span className="sp" />
      <button type="button" className="chip" onClick={toggleTheme}>
        Tema
      </button>
      <span className="av" aria-label="Perfil">
        VH
      </span>
    </header>
  );
}