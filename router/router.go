package router

import (
	"User-Management/controller"
	"User-Management/middleware"
	"User-Management/repository"
	"User-Management/service"
	"database/sql"
	"github.com/gin-gonic/gin"
)

// InitializeRoutes sets up the routes for user-related operations
func InitializeRoutes(router *gin.Engine, db *sql.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	users := router.Group("/users")
	users.Use(middleware.LoggerMiddleware())
	{
		users.POST("/", userController.CreateUser)
		users.GET("/", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}
}
