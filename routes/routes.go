package routes

import (
	"invitified-go/config"
	"invitified-go/controllers"
	"invitified-go/middlewares"
	"invitified-go/repositories"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.DB)
	tokenRepo := repositories.NewTokenRepository(config.DB)
	equipmentRepo := repositories.NewEquipmentRepository(config.DB)
	rentalRepo := repositories.NewRentalRepository(config.DB)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo, tokenRepo)
	equipmentController := controllers.NewEquipmentController(equipmentRepo)
	rentalController := controllers.NewRentalController(rentalRepo, equipmentRepo)

	// User routes
	userGroup := e.Group("/users")
	userGroup.POST("/register", userController.RegisterUser)
	userGroup.POST("/login", userController.LoginUser)

	// Protected routes
	userGroup.GET("/me", userController.GetUserProfile, middlewares.JWTMiddleware(tokenRepo))
	userGroup.DELETE("/:id", userController.DeleteUser, middlewares.JWTMiddleware(tokenRepo))

	// Equipment category routes
	categoryGroup := e.Group("/categories")
	categoryGroup.POST("", equipmentController.CreateCategory, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))
	categoryGroup.GET("/:id", equipmentController.GetCategoryByID)
	categoryGroup.GET("", equipmentController.GetAllCategories)
	categoryGroup.PUT("/:id", equipmentController.UpdateCategory, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))
	categoryGroup.DELETE("/:id", equipmentController.DeleteCategory, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))

	// Equipment routes
	equipmentGroup := e.Group("/equipment")
	equipmentGroup.POST("", equipmentController.CreateEquipment, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))
	equipmentGroup.GET("/:slug", equipmentController.GetEquipmentBySlug)
	equipmentGroup.GET("", equipmentController.GetAllEquipment)
	equipmentGroup.PUT("/:slug", equipmentController.UpdateEquipment, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))
	equipmentGroup.DELETE("/:slug", equipmentController.DeleteEquipment, middlewares.JWTMiddleware(tokenRepo), middlewares.IsAdmin(userRepo))

	// Rental routes
	rentalGroup := e.Group("/rentals")
	rentalGroup.POST("", rentalController.CreateRental, middlewares.JWTMiddleware(tokenRepo))
	rentalGroup.GET("/:id", rentalController.GetRentalByID, middlewares.JWTMiddleware(tokenRepo))
	rentalGroup.GET("", rentalController.GetAllRentals, middlewares.JWTMiddleware(tokenRepo))
	rentalGroup.PUT("/:id", rentalController.UpdateRental, middlewares.JWTMiddleware(tokenRepo))
	rentalGroup.DELETE("/:id", rentalController.DeleteRental, middlewares.JWTMiddleware(tokenRepo))
}
