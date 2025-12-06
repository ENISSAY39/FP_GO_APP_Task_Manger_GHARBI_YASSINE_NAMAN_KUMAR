package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TaskStatusTodo  = "TODO"
	TaskStatusDoing = "DOING"
	TaskStatusDone  = "DONE"

	TaskPriorityLow    = "LOW"
	TaskPriorityMedium = "MEDIUM"
	TaskPriorityHigh   = "HIGH"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProjectID   uint           `gorm:"index;not null" json:"project_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      string         `gorm:"size:20;default:TODO" json:"status"`
	Priority    string         `gorm:"size:20;default:MEDIUM" json:"priority"`
	DueDate     *time.Time     `json:"due_date"`

	CreatorID uint `gorm:"index;not null" json:"creator_id"`

	Assignees []TaskAssignee `gorm:"foreignKey:TaskID" json:"assignees"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
