package routes

import (
	"github.com/gin-gonic/gin"
)

type ICarHandler interface {
	//booking
	GetBookingsByParams(c *gin.Context)
	CreateBooking(c *gin.Context)
	UpdateBooking(c *gin.Context)
	DeleteBooking(c *gin.Context)
	//car
	GetCar(c *gin.Context)
	CreateCar(c *gin.Context)
	UpdateCar(c *gin.Context)
	DeleteCar(c *gin.Context)
}

type Router struct {
	RouterBase *gin.RouterGroup
	CarHandler ICarHandler
}

func New(r *gin.Engine, carHandler ICarHandler) *Router {
	carV1 := r.Group("/car/v1")

	return &Router{
		RouterBase: carV1,
		CarHandler: carHandler,
	}
}
