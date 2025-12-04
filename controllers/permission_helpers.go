package controllers

import (
	"errors"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

// GetMemberRole returns role string or empty + error if not member
func GetMemberRole(projectID uint, userID uint) (string, error) {
	db := initializers.DB
	var pm models.ProjectMember
	if err := db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&pm).Error; err != nil {
		return "", err
	}
	return pm.Role, nil
}

// IsProjectOwner checks if given user is the owner (Project.OwnerID) OR has a member record with RoleOwner
func IsProjectOwner(projectID uint, userID uint) (bool, error) {
	db := initializers.DB
	var project models.Project
	if err := db.First(&project, projectID).Error; err != nil {
		return false, err
	}
	if project.OwnerID != nil && *project.OwnerID == userID {
		return true, nil
	}
	// fallback: ProjectMember role
	var pm models.ProjectMember
	if err := db.Where("project_id = ? AND user_id = ? AND role = ?", projectID, userID, models.RoleOwner).First(&pm).Error; err == nil {
		return true, nil
	}
	return false, nil
}

// IsProjectMember checks if user is member (role any) of project
func IsProjectMember(projectID uint, userID uint) (bool, error) {
	db := initializers.DB
	var pm models.ProjectMember
	if err := db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&pm).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Helper to add member
func AddProjectMember(projectID uint, userID uint, role string) error {
	db := initializers.DB
	// avoid duplicates
	var existing models.ProjectMember
	if err := db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&existing).Error; err == nil {
		// already exists, update role if different
		if existing.Role != role {
			existing.Role = role
			return db.Save(&existing).Error
		}
		return nil
	}
	pm := models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	return db.Create(&pm).Error
}

func RemoveProjectMember(projectID uint, userID uint) error {
	db := initializers.DB
	return db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&models.ProjectMember{}).Error
}

// ErrNotMember sentinel
var ErrNotMember = errors.New("not a member")
