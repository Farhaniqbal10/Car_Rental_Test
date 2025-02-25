package routes

import (
	"car_rental_test/app/config"
	handlers "car_rental_test/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter inisialisasi router dan mendaftarkan routes
func SetupRouter(db *config.Database) *gin.Engine {
	r := gin.Default()

	// Buat instance CarHandler dengan Dependency Injection
	carHandler := handlers.NewCarHandler(db)

	// Daftarkan routes Cars
	RegisterCarRoutes(r, carHandler)

	return r
}

// RegisterCarRoutes mengatur semua endpoint untuk CarHandler
func RegisterCarRoutes(router *gin.Engine, carHandler *handlers.CarHandler) {
	carRoutes := router.Group("/cars")
	{
		carRoutes.GET("/", carHandler.GetCarsByParams)        // Get car by params
		carRoutes.POST("/add", carHandler.CreateCar)          // Create new car
		carRoutes.PUT("/update/:id", carHandler.UpdateCar)    // Update car by ID
		carRoutes.DELETE("/delete/:id", carHandler.DeleteCar) // Delete car by ID
	}
}
