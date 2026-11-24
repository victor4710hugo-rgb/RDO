package models

import "time"

type Project struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ClientName string    `json:"client_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WorkOrder struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Number    string    `json:"number"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Activity struct {
	ID          string    `json:"id"`
	WorkOrderID string    `json:"work_order_id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FieldReport struct {
	ID          string    `json:"id"`
	WorkOrderID string    `json:"work_order_id"`
	ReportDate  string    `json:"report_date"`
	Version     int       `json:"version"`
	Payload     string    `json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SyncRecord struct {
	ID        string    `json:"id"`
	EntityType string   `json:"entity_type"`
	EntityID  string    `json:"entity_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}	