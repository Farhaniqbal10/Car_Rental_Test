package driver

import (
	"car_rental_test/modules/car/v1/data"
	"car_rental_test/modules/car/v1/handlers"
	"car_rental_test/modules/car/v1/routes"
	"car_rental_test/modules/car/v1/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	carData    *data.Data
	carService *services.Service
	CarHandler *handlers.Handler
	carRoutes  *routes.Router
)

func initCar(r *gin.Engine, db *sqlx.DB) {
	carData = data.New(db)
	carService = services.New(carData)
	CarHandler = handlers.New(carService)

	carRoutes = routes.New(r, CarHandler)
	carRoutes.BookingRoutes()
}
