package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"` // référence vers user (owner principal)
	Owner       User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Members []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Tasks   []Task          `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}
