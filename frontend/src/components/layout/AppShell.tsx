import type { ReactNode } from "react";
import Sidebar from "./Sidebar";
import Topbar from "./Topbar";

export default function AppShell({ children }: { children: ReactNode }) {
  return (
    <div className="shell">
      <Sidebar />
      <main className="main">
        <Topbar title="RDO" subtitle="Gestão de operações de campo" />
        {children}
      </main>
      <Sidebar variant="bottom" />
    </div>
  );
}