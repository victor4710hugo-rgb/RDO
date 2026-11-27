import { useEffect, useState } from "react";
import type {
  Activity,
  FieldReport,
  FieldReportInput,
  Project,
  WorkOrder,
} from "../types";

const BASE = process.env.NEXT_PUBLIC_API_URL ?? "/api/v1";

async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${BASE}${path}`, {
    headers: { Accept: "application/json" },
  });

  if (!response.ok) {
    throw new Error(`Erro ${response.status} ao consultar ${path}`);
  }

  return (await response.json()) as T;
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const response = await fetch(`${BASE}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(`Erro ${response.status} ao enviar ${path}`);
  }

  return (await response.json()) as T;
}

export const api = {
  projects: () => get<Project[]>("/projects"),
  workOrders: () => get<WorkOrder[]>("/work-orders"),
  activities: () => get<Activity[]>("/activities"),
  fieldReports: () => get<FieldReport[]>("/field-reports"),
  saveFieldReport: (report: FieldReportInput) =>
    post<FieldReport>("/field-reports", report),
};

export function useApi<T>(loader: () => Promise<T>) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;

    loader()
      .then((result) => {
        if (alive) setData(result);
      })
      .catch((err: Error) => {
        if (alive) setError(err.message);
      })
      .finally(() => {
        if (alive) setLoading(false);
      });

    return () => {
      alive = false;
    };
  }, []);

  return { data, error, loading };
}