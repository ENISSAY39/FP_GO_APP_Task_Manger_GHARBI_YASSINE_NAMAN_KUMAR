package models

import (
	"slices"
	"time"

	"gorm.io/gorm"
)

// Vocabulaire des tâches. Ces valeurs sont exactement celles stockées en base
// et celles que le frontend envoie et affiche : voir CONTEXT.md, qui en est la
// source de vérité pour les deux côtés.
const (
	TaskStatusTodo       = "TODO"
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusDone       = "DONE"

	TaskPriorityLow    = "LOW"
	TaskPriorityMedium = "MEDIUM"
	TaskPriorityHigh   = "HIGH"
)

// AllowedTaskStatuses : les seuls statuts qu'une tâche peut porter.
var AllowedTaskStatuses = []string{TaskStatusTodo, TaskStatusInProgress, TaskStatusDone}

// AllowedTaskPriorities : les seules priorités qu'une tâche peut porter.
var AllowedTaskPriorities = []string{TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh}

// IsValidTaskStatus indique si status fait partie du vocabulaire ci-dessus.
func IsValidTaskStatus(status string) bool {
	return slices.Contains(AllowedTaskStatuses, status)
}

// IsValidTaskPriority indique si priority fait partie du vocabulaire ci-dessus.
func IsValidTaskPriority(priority string) bool {
	return slices.Contains(AllowedTaskPriorities, priority)
}

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
