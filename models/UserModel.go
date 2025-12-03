package models

import "gorm.io/gorm"

// User represents an application user
type User struct {
	gorm.Model
	Email    string         `gorm:"unique;not null" json:"email"`
	Password string         `gorm:"not null" json:"-"`
	Projects []Project      `gorm:"foreignKey:OwnerID" json:"-"`
	Tasks    []Task         `gorm:"many2many:task_assignees" json:"-"`
	Members  []ProjectMember `gorm:"foreignKey:UserID" json:"-"`
}
