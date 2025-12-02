package models

import "time"

// TaskAssignee is the explicit join table (optional).
// If you keep it, GORM will still map the many2many to the table name.
// You can add fields like AssignedAt later.
type TaskAssignee struct {
	TaskID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
}
