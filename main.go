package main

import (
	"os"
	"time"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/controllers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// load env, connect DB and run automigrations
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDataBase()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := gin.Default()

	// CORS configuration
	// IMPORTANT: when AllowCredentials is true you cannot use "*" for AllowOrigins.
	// Add whatever origin(s) you serve the frontend from (e.g. :5500 for python server or Live Server).
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:5500", // static front (python -m http.server 5500 or Live Server)
			"http://127.0.0.1:5500",
			"http://localhost:3000", // optional if you ever serve front on 3000
			"http://127.0.0.1:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve frontend static files (optional — you can also run static server separately)
	// If you want to serve the frontend from the Go server, place the 'frontend' folder next to the binary.
	router.Static("/frontend", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// ---------- API ROUTES ----------

	// auth
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	// validate returns current user — protected by middleware.RequireAuth which checks cookie JWT
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	// projects (protected endpoints require an authenticated user)
	router.POST("/projects", middleware.RequireAuth, controllers.ProjectsCreate)
	router.GET("/projects", middleware.RequireAuth, controllers.ProjectsIndex)
	router.GET("/projects/:id", middleware.RequireAuth, controllers.ProjectsShow)
	router.PUT("/projects/:id", middleware.RequireAuth, controllers.ProjectsUpdate)
	router.DELETE("/projects/:id", middleware.RequireAuth, controllers.ProjectsDelete)

	// tasks
	router.POST("/projects/:id/tasks", middleware.RequireAuth, controllers.TasksCreate)
	router.GET("/projects/:id/tasks", middleware.RequireAuth, controllers.TasksIndexForProject)
	router.GET("/tasks/:id", middleware.RequireAuth, controllers.TasksShow)
	router.PUT("/tasks/:id", middleware.RequireAuth, controllers.TasksUpdate)
	router.DELETE("/tasks/:id", middleware.RequireAuth, controllers.TasksDelete)
	router.PUT("/tasks/:id/assign", middleware.RequireAuth, controllers.TasksAssign)

	// start server
	router.Run(":" + port)
}
