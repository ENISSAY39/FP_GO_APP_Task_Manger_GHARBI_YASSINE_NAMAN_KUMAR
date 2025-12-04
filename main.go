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

	// CORS (dev-friendly)
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve frontend static files if present
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

		// Serve static assets if folder exists (common setups)
		router.Static("/assets", filepath.Join(found, "assets"))
		router.StaticFile("/", filepath.Join(found, "index.html"))

		// SPA fallback: send index.html for unknown routes
		router.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(found, "index.html"))
		})
	} else {
		log.Println("No frontend build found (tried frontend/dist, frontend/build, public). API-only mode.")
	}

	// -------------------- API GROUP --------------------
	api := router.Group("/api")
	{
		// Auth
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)

		// Validate simple token check
		api.GET("/validate", middleware.RequireAuth(), func(c *gin.Context) {
			v, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"ok": true, "userID": v})
		})

		// Logout (clears cookie if any)
		api.POST("/logout", middleware.RequireAuth(), func(c *gin.Context) {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.JSON(http.StatusOK, gin.H{"message": "logged out"})
		})

		// Projects
		api.POST("/projects", middleware.RequireAuth(), controllers.CreateProject)
		api.GET("/projects", middleware.RequireAuth(), controllers.GetMyProjects)
		api.GET("/projects/:projectId", middleware.RequireAuth(), controllers.GetProjectDetail)
		api.DELETE("/projects/:projectId", middleware.RequireAuth(), controllers.DeleteProject)

		// Members
		api.POST("/projects/:projectId/members", middleware.RequireAuth(), controllers.AddMember)
		api.DELETE("/projects/:projectId/members/:userId", middleware.RequireAuth(), controllers.RemoveMemberByParam)

		// Tasks
		api.POST("/projects/:projectId/tasks", middleware.RequireAuth(), controllers.CreateTask)
		api.GET("/projects/:projectId/tasks", middleware.RequireAuth(), controllers.GetProjectTasks)

		api.PUT("/tasks/:taskId", middleware.RequireAuth(), controllers.UpdateTask)
		api.DELETE("/tasks/:taskId", middleware.RequireAuth(), controllers.DeleteTask)

		api.PUT("/tasks/:taskId/assign", middleware.RequireAuth(), controllers.AssignTask)
		api.PUT("/tasks/:taskId/unassign", middleware.RequireAuth(), controllers.UnassignTask)
	}

	log.Printf("Starting server on :%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
