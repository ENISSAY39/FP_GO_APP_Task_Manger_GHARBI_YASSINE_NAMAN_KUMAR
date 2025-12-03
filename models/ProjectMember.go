package models



// ProjectMember links a user to a project with a role (e.g. "admin", "member")
type ProjectMember struct {
	ProjectID uint   `json:"project_id"`
	Project   Project `json:"project,omitempty"`
	UserID    uint   `json:"user_id"`
	User      User   `json:"user,omitempty"`
	Role      string `gorm:"not null" json:"role"` // e.g. "admin", "member"
}
