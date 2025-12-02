package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Signup creates a new user (hashes password)
func Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Login authenticates user and sets JWT cookie and returns if user is a global chef (admin)
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var user models.User
	if err := initializers.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// create token
	secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	// Set cookie HttpOnly
	c.SetSameSite(0)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	// check if the user is admin of at least one project (chef)
	var count int64
	initializers.DB.Model(&models.ProjectMember{}).
		Where("user_id = ? AND role = ?", user.ID, "admin").
		Count(&count)

	isChef := count > 0

	c.JSON(http.StatusOK, gin.H{"message": "logged in", "is_chef": isChef})
}

// Validate returns the authenticated user and extra info (projects/tasks where he participates)
func Validate(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user := userAny.(models.User)

	// load memberships and projects
	var memberships []models.ProjectMember
	initializers.DB.Preload("Project").Where("user_id = ?", user.ID).Find(&memberships)

	// fetch projects where user is member
	projects := make([]models.Project, 0)
	for _, m := range memberships {
		var p models.Project
		initializers.DB.Preload("Tasks").Preload("Members").Preload("Owner").First(&p, m.ProjectID)
		projects = append(projects, p)
	}

	// tasks assigned to user
	var tasks []models.Task
	initializers.DB.Preload("Project").Model(&user).Association("Tasks").Find(&tasks)

	// is chef?
	var adminCount int64
	initializers.DB.Model(&models.ProjectMember{}).
		Where("user_id = ? AND role = ?", user.ID, "admin").
		Count(&adminCount)
	isChef := adminCount > 0

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":     user.ID,
			"email":  user.Email,
			"is_chef": isChef,
		},
		"projects": projects,
		"tasks":    tasks,
	})
}
