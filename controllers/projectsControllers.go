package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// ------------------------- CREATE PROJECT -------------------------
func ProjectsCreate(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerID     uint   `json:"owner_id"` // only admin can set
		Members     []struct {
			UserID uint `json:"user_id"`
		} `json:"members"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	ownerID := user.ID
	if user.IsAdmin && body.OwnerID != 0 {
		ownerID = body.OwnerID
	}

	project := models.Project{
		Name:        body.Name,
		Description: strings.TrimSpace(body.Description),
		OwnerID:     ownerID,
	}

	if err := initializers.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	// Add members including owner
	if err := initializers.DB.Where("project_id = ? AND user_id = ?", project.ID, ownerID).
		FirstOrCreate(&models.ProjectMember{ProjectID: project.ID, UserID: ownerID, Role: "member"}).Error; err != nil {
		// not fatal for project creation, but inform
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add project owner as member"})
		return
	}

	for _, m := range body.Members {
		if m.UserID == ownerID {
			continue
		}
		if err := initializers.DB.Where("project_id = ? AND user_id = ?", project.ID, m.UserID).
			FirstOrCreate(&models.ProjectMember{ProjectID: project.ID, UserID: m.UserID, Role: "member"}).Error; err != nil {
			// skip failing member but return error (could be adjusted)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add project members"})
			return
		}
	}

	if err := initializers.DB.Preload("Owner").Preload("Members.User").Preload("Tasks").First(&project, project.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load project"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// ------------------------- LIST PROJECTS -------------------------
func ProjectsIndex(c *gin.Context) {
	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	var projects []models.Project
	var err error
	if user.IsAdmin {
		err = initializers.DB.Preload("Tasks").Preload("Owner").Preload("Members.User").Find(&projects).Error
	} else {
		err = initializers.DB.
			Preload("Tasks").
			Preload("Owner").
			Preload("Members.User").
			Joins("LEFT JOIN project_members ON project_members.project_id = projects.id").
			Where("projects.owner_id = ? OR project_members.user_id = ?", user.ID, user.ID).
			Group("projects.id").
			Find(&projects).Error
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// ------------------------- SHOW PROJECT -------------------------
func ProjectsShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var project models.Project
	if err := initializers.DB.Preload("Owner").Preload("Tasks.Assignees").
		Preload("Members.User").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	if !user.IsAdmin && user.ID != project.OwnerID {
		var count int64
		if err := initializers.DB.Model(&models.ProjectMember{}).
			Where("project_id = ? AND user_id = ?", project.ID, user.ID).
			Count(&count).Error; err != nil || count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ------------------------- UPDATE PROJECT -------------------------
func ProjectsUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
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
	user := uAny.(models.User)

	if !isAdminOrOwner(user, project.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if body.Name != "" {
		project.Name = strings.TrimSpace(body.Name)
	}
	if body.Description != "" {
		project.Description = strings.TrimSpace(body.Description)
	}

	if err := initializers.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update project"})
		return
	}
	if err := initializers.DB.Preload("Owner").Preload("Tasks").Preload("Members.User").First(&project, project.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ------------------------- DELETE PROJECT -------------------------
func ProjectsDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
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
	user := uAny.(models.User)

	if !isAdminOrOwner(user, project.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := initializers.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete project"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ------------------------- ADD MEMBERS TO PROJECT -------------------------
func ProjectsAddMembers(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var body struct {
		MemberIDs []uint `json:"member_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// verify project exists
	var project models.Project
	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	uAny, _ := c.Get("user")
	user := uAny.(models.User)

	// only admin or owner can add members
	if !isAdminOrOwner(user, project.ID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// add each member (skip duplicates)
	for _, uid := range body.MemberIDs {
		if uid == project.OwnerID {
			continue
		}
		initializers.DB.Where("project_id = ? AND user_id = ?", project.ID, uid).
			FirstOrCreate(&models.ProjectMember{ProjectID: project.ID, UserID: uid, Role: "member"})
	}

	// return updated project with members
	initializers.DB.Preload("Owner").Preload("Members.User").Preload("Tasks").First(&project, project.ID)
	c.JSON(http.StatusOK, gin.H{"project": project})
}
