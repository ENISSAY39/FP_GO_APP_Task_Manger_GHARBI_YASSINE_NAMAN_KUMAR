package controllers

import (
	"net/http"
	"strings"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// ProjectsCreate creates a project owned by the authenticated user
// Accepts JSON body:
// {
//   "name":"My project",
//   "description":"...",
//   "members": [{"user_id":2,"role":"member"}, {"user_id":3,"role":"admin"}]
// }
// Note: Owner becomes a member with role "admin" implicitly.
func ProjectsCreate(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Members     []struct {
			UserID uint   `json:"user_id"`
			Role   string `json:"role"`
		} `json:"members"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser := uAny.(models.User)

	project := models.Project{
		Name:        body.Name,
		Description: strings.TrimSpace(body.Description),
		OwnerID:     currentUser.ID,
	}

	if err := initializers.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}

	// ensure owner is admin member
	ownerMember := models.ProjectMember{
		ProjectID: project.ID,
		UserID:    currentUser.ID,
		Role:      "admin",
	}
	initializers.DB.Create(&ownerMember)

	// add other members if provided
	for _, m := range body.Members {
		// skip if user is owner
		if m.UserID == currentUser.ID {
			continue
		}
		if m.Role == "" {
			m.Role = "member"
		}
		pm := models.ProjectMember{
			ProjectID: project.ID,
			UserID:    m.UserID,
			Role:      m.Role,
		}
		initializers.DB.Create(&pm)
	}

	// reload with relations
	if err := initializers.DB.Preload("Owner").Preload("Members.User").Preload("Tasks").First(&project, project.ID).Error; err != nil {
		c.JSON(http.StatusCreated, gin.H{"project": project})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// ProjectsIndex returns all projects (preloads owner and tasks and members)
func ProjectsIndex(c *gin.Context) {
	var projects []models.Project
	if err := initializers.DB.Preload("Owner").Preload("Tasks").Preload("Members.User").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// ProjectsShow returns a single project by id (with owner and tasks and members)
func ProjectsShow(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.Preload("Owner").Preload("Tasks.Assignees").Preload("Members.User").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ProjectsUpdate and Delete similar to earlier pattern (owner only)
func ProjectsUpdate(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser := uAny.(models.User)
	if project.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	project.Name = strings.TrimSpace(body.Name)
	project.Description = strings.TrimSpace(body.Description)

	if err := initializers.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}

	if err := initializers.DB.Preload("Owner").Preload("Tasks").Preload("Members.User").First(&project, project.ID).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"project": project})
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

func ProjectsDelete(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	currentUser := uAny.(models.User)
	if project.OwnerID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if err := initializers.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete project"})
		return
	}
	c.Status(http.StatusNoContent)
}
