export type Priority = "low" | "medium" | "high" | "critical";

export type WorkforceControlMode = "aggregate" | "individual" | "both";

export type ProjectStatus =
  | "planned"
  | "in_progress"
  | "on_hold"
  | "completed"
  | "cancelled";

export type WorkOrderStatus =
  | "draft"
  | "planned"
  | "in_progress"
  | "blocked"
  | "completed"
  | "cancelled";

export type ActivityStatus =
  | "pending"
  | "in_progress"
  | "blocked"
  | "done"
  | "cancelled";

export type FieldReportStatus =
  | "draft"
  | "submitted"
  | "approved"
  | "rejected"
  | "closed";

export type SyncStatus = "pending" | "syncing" | "synced" | "error" | "conflict";

export interface Project {
  id: string;
  client_id: string;
  client_name?: string;
  name: string;
  workforce_control_mode: WorkforceControlMode;
  status: ProjectStatus;
  created_at: string;
  updated_at: string;
}

export interface WorkOrder {
  id: string;
  project_id: string;
  unit_id?: string | null;
  number: string;
  title: string;
  priority: Priority;
  status: WorkOrderStatus;
  created_at: string;
  updated_at: string;
}

export interface Activity {
  id: string;
  work_order_id: string;
  title: string;
  priority: Priority;
  status: ActivityStatus;
  planned_progress: number;
  actual_progress: number;
  created_at: string;
  updated_at: string;
}

export interface FieldReport {
  id: string;
  project_id: string;
  report_number: string;
  report_date: string;
  status: FieldReportStatus;
  version: number;
  payload?: string | null;
  created_at: string;
  updated_at: string;
}

export interface FieldReportInput {
  id: string;
  project_id: string;
  report_number: string;
  report_date: string;
  status?: FieldReportStatus;
  version: number;
  payload: string;
}

export interface SyncRecord {
  id: string;
  entity_type: string;
  entity_id: string;
  device_id?: string | null;
  operation: "create" | "update" | "delete";
  status: SyncStatus;
  record_version: number;
  attempts: number;
  last_error?: string | null;
  synced_at?: string | null;
  created_at: string;
  updated_at: string;
}