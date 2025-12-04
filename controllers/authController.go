package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
)

type signupPayload struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Signup create user
func Signup(c *gin.Context) {
	var body signupPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := initializers.DB

	// check if email already exists
	var existing models.User
	err := db.Where("email = ?", body.Email).First(&existing).Error

	switch {
	case err == nil:
		// found -> email already used
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
		return
	case errors.Is(err, gorm.ErrRecordNotFound):
		// not found -> continue
	default:
		// other DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	user := models.User{
		Name:  body.Name,
		Email: body.Email,
	}

	// use single err variable (avoid shadowing)
	if err = user.SetPassword(body.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	// minimal response without password
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}

// Login returns JWT token
func Login(c *gin.Context) {
	var body loginPayload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := initializers.DB
	var user models.User
	if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		// do not reveal whether email exists
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !user.CheckPassword(body.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signed,
		"user":  gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}
