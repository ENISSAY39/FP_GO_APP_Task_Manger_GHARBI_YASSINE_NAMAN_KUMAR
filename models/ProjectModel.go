package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleOwner  = "OWNER"
	RoleMember = "MEMBER"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`

	OwnerID *uint `gorm:"index" json:"owner_id"`
	Owner   User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"owner"`

	Members []ProjectMember `gorm:"foreignKey:ProjectID" json:"members"`
	Tasks   []Task          `gorm:"foreignKey:ProjectID" json:"tasks"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
