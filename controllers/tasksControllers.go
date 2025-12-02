package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// Create task under a project. Accepts assignee_ids []uint to assign multiple users.
func TasksCreate(c *gin.Context) {
	var body struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		Priority    int        `json:"priority"`
		AssigneeIDs []uint     `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	projectIdStr := c.Param("id")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project id"})
		return
	}

	var project models.Project
	if err := initializers.DB.First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		DueDate:     body.DueDate,
		Priority:    body.Priority,
		ProjectID:   uint(projectId),
		Status:      "todo",
	}

	if err := initializers.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}

	// assign multiple users if provided
	if len(body.AssigneeIDs) > 0 {
		var users []models.User
		initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users)
		if len(users) > 0 {
			initializers.DB.Model(&task).Association("Assignees").Replace(&users)
			initializers.DB.Preload("Assignees").First(&task, task.ID)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// TasksIndexForProject lists tasks for a project
func TasksIndexForProject(c *gin.Context) {
	projectIdStr := c.Param("id")
	var tasks []models.Task
	initializers.DB.Preload("Assignees").Where("project_id = ?", projectIdStr).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// TasksShow returns a task by id (preload assignees)
func TasksShow(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := initializers.DB.Preload("Assignees").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// TasksUpdate updates a task
func TasksUpdate(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Status      string     `json:"status"`
		DueDate     *time.Time `json:"due_date"`
		Priority    int        `json:"priority"`
		AssigneeIDs []uint     `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var task models.Task
	if err := initializers.DB.Preload("Assignees").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	// update fields if provided
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
	task.Priority = body.Priority

	if err := initializers.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	// update assignees if provided
	if body.AssigneeIDs != nil {
		var users []models.User
		initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users)
		initializers.DB.Model(&task).Association("Assignees").Replace(&users)
		initializers.DB.Preload("Assignees").First(&task, task.ID)
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// TasksDelete deletes a task
func TasksDelete(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Delete(&models.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

// TasksAssign assigns multiple users to a task
func TasksAssign(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		AssigneeIDs []uint `json:"assignee_ids"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var task models.Task
	if err := initializers.DB.Preload("Assignees").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	var users []models.User
	initializers.DB.Where("id IN ?", body.AssigneeIDs).Find(&users)
	initializers.DB.Model(&task).Association("Assignees").Replace(&users)
	initializers.DB.Preload("Assignees").First(&task, task.ID)
	c.JSON(http.StatusOK, gin.H{"task": task})
}
