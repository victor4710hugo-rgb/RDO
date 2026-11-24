"use client";

import { api, useApi } from "@/lib/api";
import { ACTIVITY_STATUS_LABEL, State } from "@/components/State";

export default function DashboardPage() {
  const projects = useApi(api.projects);
  const workOrders = useApi(api.workOrders);
  const activities = useApi(api.activities);

  const activityList = activities.data ?? [];
  const doneActivities = activityList.filter(
    (activity) => activity.status === "done",
  ).length;
  const progress =
    activityList.length > 0
      ? Math.round((doneActivities / activityList.length) * 100)
      : 0;

  return (
    <>
      <div className="grid two">
        <div className="glass stat">
          <b>{projects.data?.length ?? "—"}</b>
          <span>Projetos</span>
        </div>

        <div className="glass stat">
          <b>{workOrders.data?.length ?? "—"}</b>
          <span>Ordens de serviço</span>
        </div>

        <div className="glass stat warn">
          <b>{activityList.length - doneActivities}</b>
          <span>Atividades pendentes</span>
        </div>

        <div className="glass stat ok">
          <b>{doneActivities}</b>
          <span>Atividades concluídas</span>
        </div>
      </div>

      <div className="glass row progress-card">
        <div className="grow">
          <h3>Progresso das atividades</h3>
          <p>{progress}% concluído</p>
          <div className="progress">
            <i style={{ width: `${progress}%` }} />
          </div>
        </div>
      </div>

      <h2 className="section-title">Atividades recentes</h2>

      <State
        loading={activities.loading}
        error={activities.error}
        empty={!activities.loading && activityList.length === 0}
      />

      {activityList.slice(0, 5).map((activity) => (
        <div key={activity.id} className="glass row">
          <div className="grow">
            <h3>{activity.title}</h3>
          </div>
          <span className={`badge ${activity.status}`}>
            {ACTIVITY_STATUS_LABEL[activity.status] ?? activity.status}
          </span>
        </div>
      ))}
    </>
  );
}