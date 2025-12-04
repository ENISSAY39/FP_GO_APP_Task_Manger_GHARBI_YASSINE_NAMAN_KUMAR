package models

import "gorm.io/gorm"

// Project represents a project that contains many tasks and has many members
type Project struct {
	gorm.Model
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	OwnerID     uint            `json:"owner_id"`
	Owner       User            `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Tasks       []Task          `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
}
