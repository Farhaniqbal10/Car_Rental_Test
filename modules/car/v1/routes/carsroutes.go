package routes

func (r *Router) CarRoutes() {
	r.RouterBase.GET("/Cars", r.CarHandler.GetCar)                  // Ambil daftar mobil dengan filter
	r.RouterBase.POST("/Cars/add", r.CarHandler.CreateCar)          // Tambah mobil baru
	r.RouterBase.PUT("/Cars/update/:id", r.CarHandler.UpdateCar)    // Update mobil berdasarkan ID
	r.RouterBase.DELETE("/Cars/delete/:id", r.CarHandler.DeleteCar) // Hapus mobil berdasarkan ID
}
