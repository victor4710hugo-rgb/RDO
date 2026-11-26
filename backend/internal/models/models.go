package models

import "time"

type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID                   string    `json:"id"`
	ClientID             string    `json:"client_id"`
	ClientName           string    `json:"client_name,omitempty"`
	Name                 string    `json:"name"`
	WorkforceControlMode string    `json:"workforce_control_mode"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type WorkOrder struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UnitID    *string   `json:"unit_id,omitempty"`
	Number    string    `json:"number"`
	Title     string    `json:"title"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Activity struct {
	ID              string    `json:"id"`
	WorkOrderID     string    `json:"work_order_id"`
	Title           string    `json:"title"`
	Priority        string    `json:"priority"`
	Status          string    `json:"status"`
	PlannedProgress float64   `json:"planned_progress"`
	ActualProgress  float64   `json:"actual_progress"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type FieldReport struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	ReportNumber string    `json:"report_number"`
	ReportDate   string    `json:"report_date"`
	Status       string    `json:"status"`
	Version      int       `json:"version"`
	Payload      *string   `json:"payload,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SyncRecord struct {
	ID            string     `json:"id"`
	EntityType    string     `json:"entity_type"`
	EntityID      string     `json:"entity_id"`
	DeviceID      *string    `json:"device_id,omitempty"`
	Operation     string     `json:"operation"`
	Status        string     `json:"status"`
	RecordVersion int        `json:"record_version"`
	Attempts      int        `json:"attempts"`
	LastError     *string    `json:"last_error,omitempty"`
	SyncedAt      *time.Time `json:"synced_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Status   string `json:"status,omitempty"`
}