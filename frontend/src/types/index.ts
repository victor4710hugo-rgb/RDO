export type ActivityStatus = "pending" | "in_progress" | "done";

export type FieldReportStatus = "draft" | "submitted" | "approved" | "rejected";

export type SyncStatus = "pending" | "syncing" | "synced" | "error" | "conflict";

export interface Project {
  id: string;
  name: string;
  client_name: string;
  created_at: string;
  updated_at: string;
}

export interface WorkOrder {
  id: string;
  project_id: string;
  number: string;
  title: string;
  created_at: string;
  updated_at: string;
}

export interface Activity {
  id: string;
  work_order_id: string;
  title: string;
  status: ActivityStatus;
  created_at: string;
  updated_at: string;
}

export interface FieldReport {
  id: string;
  work_order_id: string;
  report_date: string;
  version: number;
  payload: string;
  created_at: string;
  updated_at: string;
}

export interface FieldReportInput {
  id: string;
  work_order_id: string;
  report_date: string;
  version: number;
  payload: string;
}

export interface SyncRecord {
  id: string;
  entity_type: string;
  entity_id: string;
  status: SyncStatus;
  created_at: string;
  updated_at: string;
}   