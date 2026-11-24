"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/", label: "Início" },
  { href: "/projetos", label: "Projetos" },
  { href: "/ordens-servico", label: "OS" },
  { href: "/atividades", label: "Atividades" },
];

function Links() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link) => {
        const isActive =
          link.href === "/" ? pathname === "/" : pathname.startsWith(link.href);

        return (
          <Link
            key={link.href}
            href={link.href}
            className={`nav-link${isActive ? " active" : ""}`}
          >
            <span>{link.label}</span>
          </Link>
        );
      })}
    </>
  );
}

export default function Sidebar({
  variant = "side",
}: {
  variant?: "side" | "bottom";
}) {
  if (variant === "bottom") {
    return (
      <nav className="bottom-nav glass" aria-label="Navegação principal">
        <Links />
      </nav>
    );
  }

  return (
    <aside className="sidebar glass">
      <div className="brand">RDO</div>
      <Links />
    </aside>
  );
}