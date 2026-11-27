"use client";

import { useState, type ReactNode } from "react";
import Sidebar from "./Sidebar";
import Topbar from "./Topbar";

export default function AppShell({ children }: { children: ReactNode }) {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className="app">
      <div
        className={`backdrop${menuOpen ? " open" : ""}`}
        onClick={() => setMenuOpen(false)}
      />
      <Sidebar open={menuOpen} onNavigate={() => setMenuOpen(false)} />
      <div className="col">
        <Topbar onMenuToggle={() => setMenuOpen((prev) => !prev)} />
        <main className="main">{children}</main>
      </div>
    </div>
  );
}