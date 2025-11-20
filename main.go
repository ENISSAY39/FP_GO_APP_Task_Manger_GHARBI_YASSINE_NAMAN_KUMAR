package main

import (
	"os"

	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/controllers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/middleware"
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
		port = "8080"
	}

	router := gin.Default()

	// auth
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	// projects
	router.POST("/projects", middleware.RequireAuth, controllers.ProjectsCreate)
	router.GET("/projects", controllers.ProjectsIndex)
	router.GET("/projects/:id", controllers.ProjectsShow)
	router.PUT("/projects/:id", middleware.RequireAuth, controllers.ProjectsUpdate)
	router.DELETE("/projects/:id", middleware.RequireAuth, controllers.ProjectsDelete)

	// tasks → IMPORTANT : utiliser :id pour éviter le conflit Gin
	router.POST("/projects/:id/tasks", middleware.RequireAuth, controllers.TasksCreate)
	router.GET("/projects/:id/tasks", controllers.TasksIndexForProject)

	router.GET("/tasks/:id", controllers.TasksShow)
	router.PUT("/tasks/:id", middleware.RequireAuth, controllers.TasksUpdate)
	router.DELETE("/tasks/:id", middleware.RequireAuth, controllers.TasksDelete)
	router.PUT("/tasks/:id/assign", middleware.RequireAuth, controllers.TasksAssign)

	router.Run(":" + port)
}
