// package main

// import (
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/controllers"
// 	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/initializers"
// 	"github.com/ENISSAY39/FP_GO_APP_Task_Manger_GHARBI_YASSINE_NAMAN_KUMAR/middleware"
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// func init() {
// 	initializers.LoadEnvVariables()
// 	initializers.ConnectToDB()
// 	initializers.SyncDataBase()
// }

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "3000"
// 	}

// 	router := gin.Default()

// 	// -------------------- CORS (dev-friendly) --------------------
// 	// Allow localhost / 127.0.0.1 origins for development.
// 	// NOTE: tighten this in production.
// 	router.Use(cors.New(cors.Config{
// 		AllowOriginFunc: func(origin string) bool {
// 			return strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1")
// 		},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	// -------------------- API GROUP --------------------
// 	api := router.Group("/api")
// 	{
// 		// AUTH
// 		api.POST("/signup", controllers.Signup)
// 		api.POST("/login", controllers.Login)
// 		api.GET("/validate", middleware.RequireAuth(), controllers.Validate)
// 		api.POST("/logout", middleware.RequireAuth(), controllers.Logout)

// 		// PROJECTS
// 		api.POST("/projects", middleware.RequireAuth(), controllers.ProjectsCreate)
// 		api.GET("/projects", middleware.RequireAuth(), controllers.ProjectsIndex)
// 		api.GET("/projects/:id", middleware.RequireAuth(), controllers.ProjectsShow)
// 		api.PUT("/projects/:id", middleware.RequireAuth(), controllers.ProjectsUpdate)
// 		api.DELETE("/projects/:id", middleware.RequireAuth(), controllers.ProjectsDelete)

// 		// TASKS
// 		api.POST("/projects/:id/tasks", middleware.RequireAuth(), controllers.TasksCreate)
// 		api.GET("/projects/:id/tasks", middleware.RequireAuth(), controllers.TasksIndexForProject)
// 		api.GET("/tasks/:id", middleware.RequireAuth(), controllers.TasksShow)
// 		api.PUT("/tasks/:id", middleware.RequireAuth(), controllers.TasksUpdate)
// 		api.DELETE("/tasks/:id", middleware.RequireAuth(), controllers.TasksDelete)
// 		api.PUT("/tasks/:id/assign", middleware.RequireAuth(), controllers.TasksAssign)
// 	}

// 	// -------------------- SERVE FRONTEND --------------------
// 	// Serve static files under /frontend to avoid conflict with /api routes.
// 	// Access your pages at: http://localhost:3000/frontend/pages/auth/signup.html
// 	router.Static("/frontend", "./frontend")

// 	// -------------------- RUN SERVER --------------------
// 	router.Run(":" + port)
// }


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

	// -------------------- API GROUP --------------------
	api := router.Group("/api")
	{
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)
		api.GET("/validate", middleware.RequireAuth(), controllers.Validate)
		api.POST("/logout", middleware.RequireAuth(), controllers.Logout)

		api.POST("/projects", middleware.RequireAuth(), controllers.ProjectsCreate)
		api.GET("/projects", middleware.RequireAuth(), controllers.ProjectsIndex)
		api.GET("/projects/:id", middleware.RequireAuth(), controllers.ProjectsShow)
		api.PUT("/projects/:id", middleware.RequireAuth(), controllers.ProjectsUpdate)
		api.DELETE("/projects/:id", middleware.RequireAuth(), controllers.ProjectsDelete)
		// add members to project
		api.POST("/projects/:id/members", middleware.RequireAuth(), controllers.ProjectsAddMembers)


		api.POST("/projects/:id/tasks", middleware.RequireAuth(), controllers.TasksCreate)
		api.GET("/projects/:id/tasks", middleware.RequireAuth(), controllers.TasksIndexForProject)
		api.GET("/tasks/:id", middleware.RequireAuth(), controllers.TasksShow)
		api.PUT("/tasks/:id", middleware.RequireAuth(), controllers.TasksUpdate)
		api.DELETE("/tasks/:id", middleware.RequireAuth(), controllers.TasksDelete)
		api.PUT("/tasks/:id/assign", middleware.RequireAuth(), controllers.TasksAssign)
	}

	// -------------------- STATIC FILES --------------------
	frontendDir := "./frontend"
	log.Printf("Serving frontend from: %s", frontendDir)

	// Serve static resource directories directly
	router.Static("/assets", filepath.Join(frontendDir, "assets"))
	router.Static("/js", filepath.Join(frontendDir, "js"))
	router.Static("/pages", filepath.Join(frontendDir, "pages"))

	// NoRoute: attempt to serve files under ./frontend, log attempts
	router.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		log.Printf("NoRoute hit for path: %s", reqPath)

		// If it's an API path, return 404 JSON
		if strings.HasPrefix(reqPath, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "api not found"})
			return
		}

		// If the request is root, serve index.html if exists
		if reqPath == "/" || reqPath == "" {
			indexPath := filepath.Join(frontendDir, "index.html")
			log.Printf("Trying index: %s", indexPath)
			if _, err := os.Stat(indexPath); err == nil {
				c.File(indexPath)
				return
			}
		}

		// Try mapping request path to a file inside frontendDir
		tryPath := filepath.Clean(filepath.Join(frontendDir, reqPath))
		log.Printf("Trying file: %s", tryPath)
		if fi, err := os.Stat(tryPath); err == nil && !fi.IsDir() {
			c.File(tryPath)
			return
		}

		// Also try with .html appended (common case)
		tryHtml := tryPath + ".html"
		log.Printf("Trying file with .html: %s", tryHtml)
		if fi, err := os.Stat(tryHtml); err == nil && !fi.IsDir() {
			c.File(tryHtml)
			return
		}

		// Fallback: try index.html in the requested directory
		tryIndex := filepath.Join(tryPath, "index.html")
		log.Printf("Trying index in dir: %s", tryIndex)
		if fi, err := os.Stat(tryIndex); err == nil && !fi.IsDir() {
			c.File(tryIndex)
			return
		}

		// Nothing found — return a helpful 404 that tells what was attempted
		c.String(http.StatusNotFound, "404 page not found\nTried:\n - %s\n - %s\n - %s\n", tryPath, tryHtml, tryIndex)
	})

	log.Printf("Listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
