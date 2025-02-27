package main

import (
	"car_rental_test/app/config"
	"car_rental_test/modules/v1/data"
	"car_rental_test/modules/v1/handlers"
	"car_rental_test/modules/v1/routes"
	"car_rental_test/modules/v1/services"
	"fmt"
)

func main() {
	// Inisialisasi Database
	db, err := config.NewDatabase()
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
		return
	}
	defer db.Close() // Tutup koneksi saat aplikasi berhenti

	// Inisialisasi repository
	bookingRepo := data.NewBookingRepository(db.Conn) // Sekarang ini tidak error

	// Inisialisasi service
	bookingService := services.NewBookingService(bookingRepo)

	// Inisialisasi handler
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Injeksi ke router
	r := routes.SetupRouter(bookingHandler)

	fmt.Println("Server is running on port 8081")
	r.Run(":8081")
}
