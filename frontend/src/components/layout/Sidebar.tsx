"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/", label: "Dashboard" },
  { href: "/projetos", label: "Projetos" },
  { href: "/ordens-servico", label: "Ordens de Serviço" },
  { href: "/atividades", label: "Atividades" },
  { href: "/rdo", label: "Relatórios Diários (RDO)" },
];

export default function Sidebar({
  open = false,
  onNavigate,
}: {
  open?: boolean;
  onNavigate?: () => void;
}) {
  const pathname = usePathname();

  return (
    <aside className={`sidebar${open ? " open" : ""}`}>
      <div className="brand">
        RDO
        <small>Operações de campo</small>
      </div>
      <nav aria-label="Navegação principal">
        {links.map((link) => {
          const isActive =
            link.href === "/"
              ? pathname === "/"
              : pathname.startsWith(link.href);

          return (
            <Link
              key={link.href}
              href={link.href}
              onClick={onNavigate}
              className={`nav-link${isActive ? " active" : ""}`}
            >
              {link.label}
            </Link>
          );
        })}
      </nav>
      <div className="sidebar-foot">RDO · Build Beta</div>
    </aside>
  );
}