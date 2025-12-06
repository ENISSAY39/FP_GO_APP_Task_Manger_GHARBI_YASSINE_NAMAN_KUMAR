package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/controllers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDataBase()
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// --------------------- CORS (dev-friendly) ---------------------
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Accept localhost + n'importe quel domaine (utile pour Vite, React, soutenance)
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "http://127.0.0.1") ||
				true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --------------------- FRONTEND SERVING ---------------------
	tryDirs := []string{"frontend/dist", "frontend/build", "public"}
	found := ""

	for _, d := range tryDirs {
		if stat, err := os.Stat(d); err == nil && stat.IsDir() {
			found = d
			break
		}
	}

	if found != "" {
		abs, _ := filepath.Abs(found)
		log.Printf("Serving frontend from %s\n", abs)

		// Serve static assets (React/Vite/etc.)
		router.Static("/", found)

		// SPA fallback (React Router, Vue Router, etc.)
		router.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(found, "index.html"))
		})

	} else {
		log.Println("No frontend build found (API-only mode).")
	}

	// -------------------------- API --------------------------
	api := router.Group("/api")
	{
		// Auth
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)

		// Validate token
		api.GET("/validate", middleware.RequireAuth(), func(c *gin.Context) {
			v, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"ok": true, "userID": v})
		})

			// Logout
		api.POST("/logout", middleware.RequireAuth(), controllers.Logout)

		// Projects
		api.POST("/projects", middleware.RequireAuth(), controllers.CreateProject) //marche 
		api.GET("/projects", middleware.RequireAuth(), controllers.GetMyProjects) //marche
		api.GET("/projects/:projectId", middleware.RequireAuth(), controllers.GetProjectDetail) //marche
		api.DELETE("/projects/:projectId", middleware.RequireAuth(), controllers.DeleteProject) 

		// Members
		api.POST("/projects/:projectId/members", middleware.RequireAuth(), controllers.AddMember) //marche
		api.DELETE("/projects/:projectId/members/:userId", middleware.RequireAuth(), controllers.RemoveMemberByParam) //marche
		

		// Tasks
		api.POST("/projects/:projectId/tasks", middleware.RequireAuth(), controllers.CreateTask) //marche
		api.GET("/projects/:projectId/tasks", middleware.RequireAuth(), controllers.GetProjectTasks) //marche
		api.PUT("/tasks/:taskId", middleware.RequireAuth(), controllers.UpdateTask) //marche 
		api.DELETE("/tasks/:taskId", middleware.RequireAuth(), controllers.DeleteTask) //marche 
		api.PUT("/tasks/:taskId/assign", middleware.RequireAuth(), controllers.AssignTask) //marche pas 
		api.PUT("/tasks/:taskId/unassign", middleware.RequireAuth(), controllers.UnassignTask) //marche pas
	}

	// -------------------- START SERVER --------------------
	log.Printf("Starting server on :%s\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
