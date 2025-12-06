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

//
// --------------------------- CREATE TASK ---------------------------
//

type createTaskPayload struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority"`
}

// CreateTask : n’importe quel membre du projet peut créer une tâche
func CreateTask(c *gin.Context) {

	// Récupération du projectId dans l'URL
	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)

	// Récupération de l'utilisateur authentifié (middleware RequireAuth)
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// Vérification que le user est bien membre du projet
	isMember, err := IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	// Récupération du JSON envoyé par le front
	var body createTaskPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Création de la tâche
	task := models.Task{
		Title:       body.Title,
		Description: body.Description,
		ProjectID:   projectID,
		CreatorID:   userID,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
	}

	// Priorité par défaut
	if task.Priority == "" {
		task.Priority = models.TaskPriorityMedium
	}

	if err := initializers.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

//
// --------------------------- LIST TASKS BY PROJECT ---------------------------
//

// GetProjectTasks : retourne toutes les tâches d’un projet si l’utilisateur en est membre
func GetProjectTasks(c *gin.Context) {

	pidStr := c.Param("projectId")
	pid64, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	projectID := uint(pid64)

	// User authentifié
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// Vérification : membre seulement
	isMember, err := IsProjectMember(projectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	// Chargement des tâches + assignees (avec info du User)
	var tasks []models.Task
	if err := initializers.DB.
		Where("project_id = ?", projectID).
		Preload("Assignees.User"). // <-- important pour le front : retourne les users assignés
		Find(&tasks).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

//
// --------------------------- UPDATE TASK ---------------------------
//

// Payload partiel pour update (PATCH/PUT)
type updateTaskPayload struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTask : seul le créateur OU le propriétaire du projet peut modifier
func UpdateTask(c *gin.Context) {

	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	// User authentifié
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// On récupère la tâche
	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// Vérification permission : creator OU owner du projet
	if task.CreatorID != userID {
		isOwner, err := IsProjectOwner(task.ProjectID, userID)
		if err != nil || !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator or owner can update task"})
			return
		}
	}

	// Lecture du JSON
	var body updateTaskPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mise à jour partielle
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

	// DB update
	if err := initializers.DB.Model(&task).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
		return
	}

	// Reload avec preload pour envoyer les assignees
	if err := initializers.DB.
		Preload("Assignees.User").
		First(&task, taskID).Error; err != nil {

		c.JSON(http.StatusOK, gin.H{
			"task":    task,
			"warning": "updated but failed to reload",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

//
// --------------------------- ASSIGN & UNASSIGN ---------------------------
//

// assignPayload : on assigne un user à une tâche
type assignPayload struct {
	UserID uint `json:"user_id" binding:"required"`
}

// AssignTask : un membre du projet peut assigner un autre membre
func AssignTask(c *gin.Context) {

	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	// User authentifié
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// Charger la task
	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Vérifier que celui qui assigne est membre
	isMember, err := IsProjectMember(task.ProjectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	// Lire le user cible
	var body assignPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier que le user cible est membre
	isTargetMember, err := IsProjectMember(task.ProjectID, body.UserID)
	if err != nil || !isTargetMember {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target user not a project member"})
		return
	}

	// Création du lien TaskAssignee
	ass := models.TaskAssignee{
		TaskID: taskID,
		UserID: body.UserID,
	}

	if err := initializers.DB.
		Where("task_id = ? AND user_id = ?", taskID, body.UserID).
		FirstOrCreate(&ass).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not assign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// UnassignTask : enlève un user d’une tâche
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

	// Charger la tâche
	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Vérifier que celui qui unassign est membre
	isMember, err := IsProjectMember(task.ProjectID, userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a project member"})
		return
	}

	// Lire user cible
	var body assignPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Suppression de l'assignee
	if err := initializers.DB.
		Where("task_id = ? AND user_id = ?", taskID, body.UserID).
		Delete(&models.TaskAssignee{}).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not unassign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unassigned"})
}

//
// --------------------------- DELETE TASK ---------------------------
//

// DeleteTask : seul le créateur ou l’owner du projet peut supprimer
func DeleteTask(c *gin.Context) {

	tidStr := c.Param("taskId")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	taskID := uint(tid64)

	// User authentifié
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	// Charger la tâche
	var task models.Task
	if err := initializers.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	// Permission : creator OU project owner
	if task.CreatorID != userID {
		isOwner, err := IsProjectOwner(task.ProjectID, userID)
		if err != nil || !isOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "only creator or owner can delete task"})
			return
		}
	}

	// Suppression
	if err := initializers.DB.Delete(&models.Task{}, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
