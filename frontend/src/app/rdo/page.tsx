"use client";

import { useState, type FormEvent } from "react";
import { api, useApi } from "@/lib/api";
import { State } from "@/components/State";
import { SAFETY_ITEMS, SafetyChecklist } from "@/components/rdo/SafetyChecklist";
import type { FieldReport } from "@/types";

const initialChecklist = () => SAFETY_ITEMS.map(() => false);

export default function RdoPage() {
  const projects = useApi(api.projects);
  const reports = useApi(api.fieldReports);

  const [savedList, setSavedList] = useState<FieldReport[]>([]);
  const [projectId, setProjectId] = useState("");
  const [reportNumber, setReportNumber] = useState("RDO-001");
  const [reportDate, setReportDate] = useState(
    () => new Date().toISOString().slice(0, 10),
  );
  const [weather, setWeather] = useState("Ensolarado");
  const [workforce, setWorkforce] = useState("1");
  const [notes, setNotes] = useState("");
  const [checklist, setChecklist] = useState<boolean[]>(initialChecklist);
  const [saving, setSaving] = useState(false);
  const [feedback, setFeedback] = useState<string | null>(null);

  const allReports = [...savedList, ...(reports.data ?? [])];
  const selectedProjectId = projectId || projects.data?.[0]?.id || "";

  function toggleItem(index: number) {
    setChecklist((prev) => prev.map((v, i) => (i === index ? !v : v)));
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!selectedProjectId) return;

    setSaving(true);
    setFeedback(null);
    try {
      const payload = JSON.stringify({
        weather,
        workforce: Number(workforce) || 0,
        notes,
        checklist: SAFETY_ITEMS.map((item, i) => ({ item, ok: checklist[i] })),
      });

      const saved = await api.saveFieldReport({
        id: crypto.randomUUID(),
        project_id: selectedProjectId,
        report_number: reportNumber.trim(),
        report_date: reportDate,
        status: "submitted",
        version: 1,
        payload,
      });

      setSavedList((prev) => [saved, ...prev]);
      setNotes("");
      setChecklist(initialChecklist());
      setFeedback(`RDO ${saved.report_number} salvo com sucesso.`);
    } catch (err) {
      setFeedback(err instanceof Error ? err.message : "Erro ao salvar RDO.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="glass" style={{ padding: 16, marginBottom: 16 }}>
        <h2 className="section-title" style={{ marginTop: 0 }}>Novo RDO</h2>

        <div style={{ display: "grid", gap: 10 }}>
          <label>
            Projeto
            <select
              value={selectedProjectId}
              onChange={(e) => setProjectId(e.target.value)}
              required
              style={{ width: "100%", padding: 8, marginTop: 4 }}
            >
              {(projects.data ?? []).map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </select>
          </label>

          <div className="grid two">
            <label>
              Número do RDO
              <input
                value={reportNumber}
                onChange={(e) => setReportNumber(e.target.value)}
                required
                style={{ width: "100%", padding: 8, marginTop: 4 }}
              />
            </label>
            <label>
              Data
              <input
                type="date"
                value={reportDate}
                onChange={(e) => setReportDate(e.target.value)}
                required
                style={{ width: "100%", padding: 8, marginTop: 4 }}
              />
            </label>
          </div>

          <div className="grid two">
            <label>
              Clima
              <select
                value={weather}
                onChange={(e) => setWeather(e.target.value)}
                style={{ width: "100%", padding: 8, marginTop: 4 }}
              >
                <option value="Ensolarado">Ensolarado</option>
                <option value="Nublado">Nublado</option>
                <option value="Chuvoso">Chuvoso</option>
              </select>
            </label>
            <label>
              Efetivo em campo
              <input
                type="number"
                min="0"
                value={workforce}
                onChange={(e) => setWorkforce(e.target.value)}
                style={{ width: "100%", padding: 8, marginTop: 4 }}
              />
            </label>
          </div>

          <label>
            Ocorrências / Observações
            <textarea
              rows={3}
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Resumo do turno, desvios ou ocorrências..."
              style={{ width: "100%", padding: 8, marginTop: 4 }}
            />
          </label>
        </div>

        <SafetyChecklist checked={checklist} onToggle={toggleItem} />

        <button type="submit" disabled={saving || !selectedProjectId} style={{ padding: "10px 16px" }}>
          {saving ? "Salvando..." : "Salvar RDO"}
        </button>

        {feedback && <p style={{ marginTop: 10 }}>{feedback}</p>}
      </form>

      <h2 className="section-title">RDOs emitidos</h2>
      <State
        loading={reports.loading}
        error={reports.error}
        empty={!reports.loading && !reports.error && allReports.length === 0}
      />

      {allReports.map((rdo) => (
        <div key={rdo.id} className="glass row">
          <div className="ico-box">RDO</div>
          <div className="grow">
            <h3>{rdo.report_number}</h3>
            <p>Data: {rdo.report_date} · v{rdo.version}</p>
          </div>
          <span className={`badge ${rdo.status}`}>{rdo.status}</span>
        </div>
      ))}
    </>
  );
}