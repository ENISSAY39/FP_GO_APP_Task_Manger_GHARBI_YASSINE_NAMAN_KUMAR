package models

import "time"

// TaskAssignee is the explicit join table for many2many relation between tasks and users.
type TaskAssignee struct {
	TaskID     uint      `gorm:"primaryKey" json:"task_id"`
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
}
