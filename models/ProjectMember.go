package models

import "gorm.io/gorm"

// ProjectMember links a user to a project with a role (e.g. "manager", "member")
type ProjectMember struct {
	gorm.Model
	ProjectID uint    `gorm:"not null;index" json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	UserID    uint    `gorm:"not null;index" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string  `gorm:"not null;default:'member'" json:"role"` // "manager" or "member"
}
