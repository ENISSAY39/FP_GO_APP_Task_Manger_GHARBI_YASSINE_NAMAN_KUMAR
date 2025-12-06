package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

type createProjectPayload struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateProject: authenticated user becomes Owner and member(RoleOwner)
func CreateProject(c *gin.Context) {
	var body createProjectPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	db := initializers.DB
	project := models.Project{
		Name:        body.Name,
		Description: body.Description,
		OwnerID:     &userID,
	}
	if err := db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create project"})
		return
	}

	// add ProjectMember as OWNER
	if err := AddProjectMember(project.ID, userID, models.RoleOwner); err != nil {
		// best-effort: log, but don't fail creation
		c.JSON(http.StatusCreated, gin.H{"project": project, "warning": "created but could not add member"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// GetMyProjects returns projects where user is a member (preload members and tasks count)
func GetMyProjects(c *gin.Context) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}
	db := initializers.DB

	var memberships []models.ProjectMember
	if err := db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	projectIDs := make([]uint, 0, len(memberships))
	for _, m := range memberships {
		projectIDs = append(projectIDs, m.ProjectID)
	}

	var projects []models.Project
	if len(projectIDs) > 0 {
		if err := db.Where("id IN ?", projectIDs).Preload("Members").Preload("Tasks").Find(&projects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load projects"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetProjectDetail returns project if user is member
func GetProjectDetail(c *gin.Context) {
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// check membership using explicit error handling
	isMember, err := IsProjectMember(projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isMember {
		// defensive: treat as forbidden
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	db := initializers.DB
	var project models.Project
	if err := db.Preload("Members").Preload("Tasks.Assignees").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// AddMember: only OWNER can add
type addMemberPayload struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role"` // optional: OWNER or MEMBER
}

func AddMember(c *gin.Context) {
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	owner, err := IsProjectOwner(projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !owner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can add members"})
		return
	}

	var body addMemberPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role := body.Role
	if role != models.RoleOwner {
		role = models.RoleMember
	}

	if err := AddProjectMember(projectID, body.UserID, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add member"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

// RemoveMember: only OWNER can remove (cannot remove self if last owner)
type removeMemberPayload struct {
	UserID uint `json:"user_id" binding:"required"`
}

// removeMemberHandlerCommon centralises the actual deletion logic and checks
func removeMemberHandlerCommon(c *gin.Context, projectID uint, targetUserID uint) {
	// caller check
	callerID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// only owner allowed
	isOwner, err := IsProjectOwner(projectID, callerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can remove members"})
		return
	}

	// don't allow removing the owner
	var proj models.Project
	if err := initializers.DB.First(&proj, projectID).Error; err == nil {
		if proj.OwnerID != nil && *proj.OwnerID == targetUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot remove project owner"})
			return
		}
	}

	// perform deletion by matching both project_id and user_id
	res := initializers.DB.
		Where("project_id = ? AND user_id = ?", projectID, targetUserID).
		Delete(&models.ProjectMember{})

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove member"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found for that project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

func RemoveMember(c *gin.Context) {
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)

	var body removeMemberPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	removeMemberHandlerCommon(c, projectID, body.UserID)
}

// DeleteProject: only owner can delete
func DeleteProject(c *gin.Context) {
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	isOwner, err := IsProjectOwner(projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete project"})
		return
	}

	if err := initializers.DB.Delete(&models.Project{}, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete project"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

// helper to read userID from context (expects uint)
func getUserIDFromCtx(c *gin.Context) (uint, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	switch id := v.(type) {
	case uint:
		return id, true
	case int:
		return uint(id), true
	case int64:
		return uint(id), true
	default:
		return 0, false
	}
}

// RemoveMemberByParam handles DELETE /projects/:projectId/members/:userId
func RemoveMemberByParam(c *gin.Context) {
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)

	uidStr := c.Param("userId")
	uid64, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	targetUserID := uint(uid64)

	removeMemberHandlerCommon(c, projectID, targetUserID)
}
