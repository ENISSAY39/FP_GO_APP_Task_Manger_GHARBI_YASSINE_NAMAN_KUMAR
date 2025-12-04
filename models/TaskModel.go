package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TaskStatusTodo     = "TODO"
	TaskStatusDoing    = "DOING"
	TaskStatusDone     = "DONE"
	TaskStatusBlocked  = "BLOCKED"
	TaskPriorityLow    = "LOW"
	TaskPriorityMedium = "MEDIUM"
	TaskPriorityHigh   = "HIGH"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	ProjectID   uint           `gorm:"index;not null" json:"project_id"`
	CreatorID   uint           `gorm:"index;not null" json:"creator_id"` // qui a créé la tâche
	Status      string         `gorm:"size:30;not null;default:'TODO'" json:"status"`
	Priority    string         `gorm:"size:20;not null;default:'MEDIUM'" json:"priority"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Project    Project        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Creator    User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Assignees  []TaskAssignee `gorm:"foreignKey:TaskID" json:"assignees,omitempty"`
}
