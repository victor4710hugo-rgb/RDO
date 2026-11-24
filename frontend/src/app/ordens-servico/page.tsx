"use client";

import { api, useApi } from "@/lib/api";
import { State } from "@/components/State";

export default function WorkOrdersPage() {
  const workOrders = useApi(api.workOrders);
  const activities = useApi(api.activities);

  return (
    <>
      <State
        loading={workOrders.loading}
        error={workOrders.error}
        empty={
          !workOrders.loading && !workOrders.error && !workOrders.data?.length
        }
      />

      {(workOrders.data ?? []).map((workOrder) => {
        const relatedActivities = (activities.data ?? []).filter(
          (activity) => activity.work_order_id === workOrder.id,
        );

        const doneActivities = relatedActivities.filter(
          (activity) => activity.status === "done",
        ).length;

        const progress =
          relatedActivities.length > 0
            ? Math.round((doneActivities / relatedActivities.length) * 100)
            : 0;

        return (
          <div key={workOrder.id} className="glass row">
            <div className="ico-box">OS</div>
            <div className="grow">
              <h3>
                {workOrder.number} · {workOrder.title}
              </h3>
              <p>
                {doneActivities} de {relatedActivities.length} atividades
                concluídas
              </p>
              <div className="progress">
                <i style={{ width: `${progress}%` }} />
              </div>
            </div>
          </div>
        );
      })}
    </>
  );
}