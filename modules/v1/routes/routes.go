package routes

import (
	"car_rental_test/modules/v1/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter inisialisasi router dengan handler
func SetupRouter(bookingHandler *handlers.BookingHandler) *gin.Engine {
	r := gin.Default()

	// Daftarkan routes
	RegisterBookingRoutes(r, bookingHandler)

	return r
}

// RegisterBookingRoutes mengatur semua endpoint terkait pemesanan (bookings)
func RegisterBookingRoutes(router *gin.Engine, bookingHandler *handlers.BookingHandler) {
	bookingRoutes := router.Group("/bookings")
	{
		bookingRoutes.GET("/", bookingHandler.GetBookingsByParams)        // Ambil daftar pemesanan berdasarkan filter
		bookingRoutes.POST("/add", bookingHandler.CreateBooking)          // Tambah pemesanan baru
		bookingRoutes.PUT("/update/:id", bookingHandler.UpdateBooking)    // Perbarui data pemesanan berdasarkan ID
		bookingRoutes.DELETE("/delete/:id", bookingHandler.DeleteBooking) // Hapus pemesanan berdasarkan ID
	}
}
