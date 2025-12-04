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

// ------------------------- SIGNUP -------------------------
func Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: string(hash),
		IsAdmin:  false,
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// ------------------------- LOGIN -------------------------
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "dev-secret"
	}

	exp := time.Now().Add(24 * time.Hour * 30).Unix()
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	// Cookie is HttpOnly and secure flag depends on env; keep SameSite lax
	c.SetSameSite(http.SameSiteLaxMode)
	// 3600*24*30 seconds = 30 days
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":  "logged in",
		"is_admin": user.IsAdmin,
	})
}

// ------------------------- LOGOUT -------------------------
func Logout(c *gin.Context) {
	// delete cookie
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ------------------------- VALIDATE -------------------------
func Validate(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := userAny.(models.User)

	var memberships []models.ProjectMember
	if err := initializers.DB.Preload("Project").Where("user_id = ?", user.ID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch memberships"})
		return
	}

	var projects []models.Project
	for _, m := range memberships {
		var p models.Project
		if err := initializers.DB.
			Preload("Tasks").
			Preload("Members.User").
			Preload("Owner").
			First(&p, m.ProjectID).Error; err == nil {
			projects = append(projects, p)
		}
	}

	var tasks []models.Task
	// load tasks assigned to user via association
	if err := initializers.DB.Preload("Project").Model(&user).Association("Tasks").Find(&tasks); err != nil {
		// GORM's Association.Find returns error type, so handle it (if any)
		// If there's an error, continue with empty tasks but notify in response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
		"projects": projects,
		"tasks":    tasks,
	})
}
