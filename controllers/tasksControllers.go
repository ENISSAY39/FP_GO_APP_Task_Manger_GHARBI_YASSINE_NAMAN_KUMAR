package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
)

// TasksCreate creates a task under a project
func TasksCreate(c *gin.Context) {
	var body struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		DueDate     *time.Time `json:"due_date"`
		Priority    int        `json:"priority"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// IMPORTANT : remplacer projectId → id
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

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// TasksIndexForProject lists tasks for a project
func TasksIndexForProject(c *gin.Context) {
	// IMPORTANT : remplacer projectId → id
	projectIdStr := c.Param("id")

	var tasks []models.Task
	initializers.DB.Where("project_id = ?", projectIdStr).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// TasksShow returns a task by id (preload assignee)
func TasksShow(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := initializers.DB.Preload("Assignee").First(&task, id).Error; err != nil {
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
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var task models.Task
	if err := initializers.DB.First(&task, id).Error; err != nil {
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
	initializers.DB.Save(&task)
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

// TasksAssign assigns a task to a user
func TasksAssign(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		AssigneeID uint `json:"assignee_id"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var task models.Task
	if err := initializers.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	// verify assignee exists
	var user models.User
	if err := initializers.DB.First(&user, body.AssigneeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignee not found"})
		return
	}
	task.AssigneeID = &body.AssigneeID
	initializers.DB.Save(&task)
	c.JSON(http.StatusOK, gin.H{"task": task})
}
