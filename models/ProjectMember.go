package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"index;not null" json:"project_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Role      string         `gorm:"size:20;not null" json:"role"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
