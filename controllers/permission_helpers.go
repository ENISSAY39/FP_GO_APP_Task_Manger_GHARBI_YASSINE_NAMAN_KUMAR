package controllers

import (
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

// isAdminOrOwner returns true if user is global admin OR owner of the project.
func isAdminOrOwner(user models.User, projectID uint) bool {
	if user.IsAdmin {
		return true
	}
	var project models.Project
	if err := initializers.DB.Select("id, owner_id").First(&project, projectID).Error; err == nil {
		if project.OwnerID == user.ID {
			return true
		}
	}
	return false
}
