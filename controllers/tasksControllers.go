package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// ------------------------- CREATE TASK -------------------------
func TasksCreate(c *gin.Context) {
	var body struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		Priority    int        `json:"priority"`
		AssigneeIDs []uint     `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	projectIdStr := c.Param("id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var project models.Project
	if err := initializers.DB.First(&project, projectId).Error; err != nil {
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

	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
		Priority:    body.Priority,
		ProjectID:   project.ID,
		Status:      "todo",
	}
	if err := initializers.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	// assign
	if len(body.AssigneeIDs) > 0 {
		var users []models.User
		if err := initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assignees"})
			return
		}
		if err := initializers.DB.Model(&task).Association("Assignees").Replace(&users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign users"})
			return
		}
		if err := initializers.DB.Preload("Assignees").First(&task, task.ID).Error; err != nil {
			// rare, but handle
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load created task"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// ------------------------- LIST TASKS FOR PROJECT -------------------------
func TasksIndexForProject(c *gin.Context) {
	projectIdStr := c.Param("id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var tasks []models.Task
	if err := initializers.DB.Preload("Assignees").Where("project_id = ?", projectId).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// ------------------------- SHOW TASK -------------------------
func TasksShow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var task models.Task
	if err := initializers.DB.Preload("Assignees").Preload("Project").
		First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	// admin always allowed
	if user.IsAdmin {
		c.JSON(http.StatusOK, gin.H{"task": task})
		return
	}

	// owner allowed
	if task.Project.OwnerID == user.ID {
		c.JSON(http.StatusOK, gin.H{"task": task})
		return
	}

	// assignee allowed
	var count int64
	if err := initializers.DB.Table("task_assignees").
		Where("task_id = ? AND user_id = ?", task.ID, user.ID).
		Count(&count).Error; err == nil && count > 0 {
		c.JSON(http.StatusOK, gin.H{"task": task})
		return
	}

	// member of project allowed
	if err := initializers.DB.Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", task.ProjectID, user.ID).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// ------------------------- UPDATE TASK -------------------------
func TasksUpdate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var body struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"`
		Priority    int        `json:"priority"`
		AssigneeIDs []uint     `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var task models.Task
	if err := initializers.DB.Preload("Assignees").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	if !isAdminOrOwner(user, task.ProjectID) {
		// Only assignee can update status
		if body.Status == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var count int64
		if err := initializers.DB.Table("task_assignees").
			Where("task_id = ? AND user_id = ?", task.ID, user.ID).
			Count(&count).Error; err != nil || count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		task.Status = body.Status
	} else {
		// Admin/Owner can edit everything
		if body.Title != "" {
			task.Title = body.Title
		}
		if body.Description != "" {
			task.Description = body.Description
		}
		if body.Status != "" {
			task.Status = body.Status
		}
		if body.DueDate != nil {
			task.DueDate = body.DueDate
		}
		// priority may be zero intentionally; we still set it
		task.Priority = body.Priority
	}

	// Save
	if err := initializers.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	// Update assignees if admin/owner
	if body.AssigneeIDs != nil && isAdminOrOwner(user, task.ProjectID) {
		var users []models.User
		if err := initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assignees"})
			return
		}
		if err := initializers.DB.Model(&task).Association("Assignees").Replace(&users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update assignees"})
			return
		}
		if err := initializers.DB.Preload("Assignees").First(&task, task.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load task"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// ------------------------- DELETE TASK -------------------------
func TasksDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var task models.Task
	if err := initializers.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	if !isAdminOrOwner(user, task.ProjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := initializers.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ------------------------- ASSIGN TASK -------------------------
func TasksAssign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var body struct {
		AssigneeIDs []uint `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var task models.Task
	if err := initializers.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	uAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := uAny.(models.User)

	if !isAdminOrOwner(user, task.ProjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var users []models.User
	if err := initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assignees"})
		return
	}
	if err := initializers.DB.Model(&task).Association("Assignees").Replace(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign users"})
		return
	}
	if err := initializers.DB.Preload("Assignees").First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}
