package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// --- 1) Récupérer le cookie ---
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			c.Abort()
			return
		}

		// --- 2) Récupérer la clé secrète ---
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			secret = "dev-secret"
		}

		// --- 3) Décoder le token JWT ---
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// on vérifie qu’on utilise bien HS256
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// --- 4) Extraire les claims ---
		claims := token.Claims.(jwt.MapClaims)

		// Vérifier expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.Abort()
			return
		}

		// Récupérer l’ID de l'utilisateur (sub)
		userID := uint(claims["sub"].(float64))

		// --- 5) Charger l'utilisateur depuis la DB ---
		var user models.User
		if err := initializers.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// --- 6) Attacher le user dans le contexte ---
		c.Set("user", user)

		// --- 7) Continuer la requête ---
		c.Next()
	}
}
