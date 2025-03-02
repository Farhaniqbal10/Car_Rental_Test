package routes

func (r *Router) BookingRoutes() {
	r.RouterBase.GET("/Bookings", r.CarHandler.GetBookingsByParams)
	r.RouterBase.POST("/Bookings/add", r.CarHandler.CreateBooking)
	r.RouterBase.PUT("/Bookings/update/:id", r.CarHandler.UpdateBooking)
	r.RouterBase.DELETE("/Bookings/delete/:id", r.CarHandler.DeleteBooking)
}
