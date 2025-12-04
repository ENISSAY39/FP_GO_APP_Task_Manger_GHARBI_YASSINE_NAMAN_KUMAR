package models

import (
	"time"

	"gorm.io/gorm"
)

// Role values: OWNER or MEMBER
const (
	RoleOwner  = "OWNER"
	RoleMember = "MEMBER"
)

type ProjectMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProjectID uint           `gorm:"index;not null" json:"project_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Role      string         `gorm:"size:20;not null;default:'MEMBER'" json:"role"` // OWNER or MEMBER
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User    User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Project Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
