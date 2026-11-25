"use client";

import { api, useApi } from "@/lib/api";
import { ACTIVITY_STATUS_LABEL, State } from "@/components/State";

export default function ActivitiesPage() {
  const { data, loading, error } = useApi(api.activities);

  return (
    <>
      <State
        loading={loading}
        error={error}
        empty={!loading && !error && !data?.length}
      />

      {(data ?? []).map((activity) => (
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