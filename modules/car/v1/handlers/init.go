package handlers

import (
	"car_rental_test/modules/car/v1/models"
	"context"
)

type ICarSvc interface {
	//booking
	GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.Booking, error)
	CreateBooking(ctx context.Context, input models.Booking) error
	UpdateBooking(ctx context.Context, id int64, input models.Booking) error
	DeleteBooking(ctx context.Context, id int64) (int, error)
	//car
	GetCar(ctx context.Context, query models.CarQueryParams) ([]models.Car, error)
	CreateCar(ctx context.Context, input models.Car) (int64, error)
	UpdateCar(ctx context.Context, id int64, input models.Car) error
	DeleteCar(ctx context.Context, id int64) (int, error)
}

type Handler struct {
	carSvc ICarSvc
}

func New(carSvc ICarSvc) *Handler {
	return &Handler{
		carSvc: carSvc,
	}
}
