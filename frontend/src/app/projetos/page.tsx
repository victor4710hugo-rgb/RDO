"use client";

import { useState } from "react";
import { api, useApi } from "@/lib/api";
import { State } from "@/components/State";

export default function ProjectsPage() {
  const { data, loading, error } = useApi(api.projects);
  const [search, setSearch] = useState("");

  const projects = (data ?? []).filter((project) =>
    `${project.name} ${project.client_name}`
      .toLowerCase()
      .includes(search.toLowerCase()),
  );

  return (
    <>
      <label className="glass search">
        <input
          placeholder="Pesquisar projeto ou cliente…"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
        />
      </label>

      <State
        loading={loading}
        error={error}
        empty={!loading && !error && projects.length === 0}
      />

      {projects.map((project) => (
        <div key={project.id} className="glass row">
          <div className="ico-box">P</div>
          <div className="grow">
            <h3>{project.name}</h3>
            <p>Cliente: {project.client_name}</p>
          </div>
        </div>
      ))}
    </>
  );
}