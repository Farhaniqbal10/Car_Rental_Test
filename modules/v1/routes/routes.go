package routes

import (
	"car_rental_test/app/config"
	"car_rental_test/modules/v1/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter inisialisasi router dan mendaftarkan routes
func SetupRouter(db *config.Database) *gin.Engine {
	r := gin.Default()

	// Buat instance handler dengan Dependency Injection
	carHandler := handlers.NewCarHandler(db)
	customerHandler := handlers.NewCustomerHandler(db)
	bookingHandler := handlers.NewBookingHandler(db) // Tambahkan handler booking

	// Daftarkan routes
	RegisterCarRoutes(r, carHandler)
	RegisterCustomerRoutes(r, customerHandler)
	RegisterBookingRoutes(r, bookingHandler) // Daftarkan routes booking

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

// RegisterCustomerRoutes mengatur semua endpoint untuk CustomerHandler
func RegisterCustomerRoutes(router *gin.Engine, customerHandler *handlers.CustomerHandler) {
	customerRoutes := router.Group("/customers")
	{
		customerRoutes.GET("/", customerHandler.GetCustomersByParams)        // Get customers by params
		customerRoutes.POST("/add", customerHandler.CreateCustomer)          // Create new customer
		customerRoutes.PUT("/update/:id", customerHandler.UpdateCustomer)    // Update customer by ID
		customerRoutes.DELETE("/delete/:id", customerHandler.DeleteCustomer) // Delete customer by ID
	}
}

// RegisterBookingRoutes mengatur semua endpoint untuk BookingHandler
func RegisterBookingRoutes(router *gin.Engine, bookingHandler *handlers.BookingHandler) {
	bookingRoutes := router.Group("/bookings")
	{
		bookingRoutes.GET("/", bookingHandler.GetBookingsByParams)        // Get bookings by params
		bookingRoutes.POST("/add", bookingHandler.CreateBooking)          // Create new booking
		bookingRoutes.PUT("/update/:id", bookingHandler.UpdateBooking)    // Update booking by ID
		bookingRoutes.DELETE("/delete/:id", bookingHandler.DeleteBooking) // Delete booking by ID
	}
}
