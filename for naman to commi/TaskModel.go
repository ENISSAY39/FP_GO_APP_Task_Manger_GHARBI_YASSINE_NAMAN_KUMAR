package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in a project. A task can be assigned to multiple users.
type Task struct {
	gorm.Model
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Status      string     `gorm:"default:todo" json:"status"` // todo, doing, done...
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ProjectID   uint       `json:"project_id"`
	Project     Project    `json:"project,omitempty"`
	Assignees   []User     `gorm:"many2many:task_assignees" json:"assignees,omitempty"`
}
