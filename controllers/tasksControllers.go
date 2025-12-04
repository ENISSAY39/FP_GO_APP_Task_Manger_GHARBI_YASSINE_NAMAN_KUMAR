package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

type createTaskPayload struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority"`
}

// CreateTask: any project member can create a task
func CreateTask(c *gin.Context) {
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

	// must be member
	isMember, err := IsProjectMember(projectID, userID)
	if err != nil {
		// DB error vs not a member
		if errors.Is(err, gorm.ErrRecordNotFound) || !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	var body createTaskPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		ProjectID:   projectID,
		CreatorID:   userID,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
	}
	if task.Priority == "" {
		task.Priority = models.TaskPriorityMedium
	}
	db := initializers.DB
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}

	// return the created task (with ID, timestamps)
	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// GetProjectTasks returns tasks in project for members
func GetProjectTasks(c *gin.Context) {
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

	isMember, err := IsProjectMember(projectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	var tasks []models.Task
	if err := initializers.DB.Where("project_id = ?", projectID).Preload("Assignees").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateTask: creator or project owner can update
type updateTaskPayload struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

func UpdateTask(c *gin.Context) {
	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// check permission: creator or project owner
	if task.CreatorID != userID {
		isOwner, err := IsProjectOwner(task.ProjectID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator or owner can update task"})
			return
		}
	}

	var body updateTaskPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := map[string]interface{}{}
	if body.Title != nil {
		updated["title"] = *body.Title
	}
	if body.Description != nil {
		updated["description"] = *body.Description
	}
	if body.Status != nil {
		updated["status"] = *body.Status
	}
	if body.Priority != nil {
		updated["priority"] = *body.Priority
	}
	if body.DueDate != nil {
		updated["due_date"] = body.DueDate
	}

	if len(updated) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := initializers.DB.Model(&task).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
		return
	}

	// reload task to return current values
	if err := initializers.DB.Preload("Assignees").First(&task, taskID).Error; err != nil {
		// even if reload fails, return success (task updated), but inform
		c.JSON(http.StatusOK, gin.H{"task": task, "warning": "updated but failed to reload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// Assign/Unassign
type assignPayload struct {
	UserID uint `json:"user_id" binding:"required"`
}

func AssignTask(c *gin.Context) {
	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// load task
	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// only project members can assign (and assign target must be a member)
	isMember, err := IsProjectMember(task.ProjectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	var body assignPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check target is member
	targetIsMember, err := IsProjectMember(task.ProjectID, body.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || !targetIsMember {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target user not a project member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !targetIsMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target user not a project member"})
		return
	}

	ass := models.TaskAssignee{
		TaskID: taskID,
		UserID: body.UserID,
	}
	if err := initializers.DB.Where("task_id = ? AND user_id = ?", taskID, body.UserID).FirstOrCreate(&ass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not assign"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func UnassignTask(c *gin.Context) {
	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// only project members can unassign
	isMember, err := IsProjectMember(task.ProjectID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	var body assignPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Where("task_id = ? AND user_id = ?", taskID, body.UserID).Delete(&models.TaskAssignee{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not unassign"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unassigned"})
}

// DeleteTask: creator or project owner can delete
func DeleteTask(c *gin.Context) {
	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if task.CreatorID != userID {
		isOwner, err := IsProjectOwner(task.ProjectID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator or owner can delete task"})
			return
		}
	}

	if err := initializers.DB.Delete(&models.Task{}, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
