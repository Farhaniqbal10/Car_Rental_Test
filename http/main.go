package main

import (
	"car_rental_test/app/config"
	"car_rental_test/modules/v1/routes"
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

	// Injeksi database ke router
	r := routes.SetupRouter(db)

	fmt.Println("Server is running on port 8081")
	r.Run(":8081")
}
