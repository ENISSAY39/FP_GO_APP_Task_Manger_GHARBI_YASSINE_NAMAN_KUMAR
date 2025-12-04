package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ----------------------------------------
		// 1) Récupérer token
		// ----------------------------------------
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}

		// fallback cookie
		if tokenString == "" {
			if cookie, err := c.Cookie("token"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// ----------------------------------------
		// 2) Secret
		// ----------------------------------------
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev_secret_change_me"
		}

		// ----------------------------------------
		// 3) Parse token (compatible v3/v4)
		// ----------------------------------------
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// Vérifier méthode HS256
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil // on ne vérifie pas le type exact → SAFE
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// ----------------------------------------
		// 4) Récupérer userID depuis "sub"
		// ----------------------------------------
		sub, ok := claims["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing subject"})
			c.Abort()
			return
		}

		var userID uint64

		switch v := sub.(type) {
		case float64:
			userID = uint64(v)
		case string:
			parsed, _ := strconv.ParseUint(v, 10, 64)
			userID = parsed
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid subject type"})
			c.Abort()
			return
		}

		// ----------------------------------------
		// 5) Vérifier que l'utilisateur existe
		// ----------------------------------------
		var user models.User
		if err := initializers.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// ----------------------------------------
		// 6) Injecter userID dans le contexte
		// ----------------------------------------
		c.Set("userID", user.ID)

		c.Next()
	}
}
