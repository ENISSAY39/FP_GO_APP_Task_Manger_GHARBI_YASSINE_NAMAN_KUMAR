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
		port = "3000"
	}

	router := gin.Default()

	// auth
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	// validate / dashboard (uses RequireAuth)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	// projects
	router.POST("/projects", middleware.RequireAuth, controllers.ProjectsCreate)
	router.GET("/projects", middleware.RequireAuth, controllers.ProjectsIndex)
	router.GET("/projects/:id", middleware.RequireAuth, controllers.ProjectsShow)
	router.PUT("/projects/:id", middleware.RequireAuth, controllers.ProjectsUpdate)
	router.DELETE("/projects/:id", middleware.RequireAuth, controllers.ProjectsDelete)

	// tasks (project ID param is :id to match project routes)
	router.POST("/projects/:id/tasks", middleware.RequireAuth, controllers.TasksCreate)
	router.GET("/projects/:id/tasks", middleware.RequireAuth, controllers.TasksIndexForProject)
	router.GET("/tasks/:id", middleware.RequireAuth, controllers.TasksShow)
	router.PUT("/tasks/:id", middleware.RequireAuth, controllers.TasksUpdate)
	router.DELETE("/tasks/:id", middleware.RequireAuth, controllers.TasksDelete)
	router.PUT("/tasks/:id/assign", middleware.RequireAuth, controllers.TasksAssign) // replace assignees

	router.Run(":" + port)
	router.Static("/frontend", "./frontend")

router.GET("/", func(c *gin.Context) {
    c.File("./frontend/index.html")
})

}
