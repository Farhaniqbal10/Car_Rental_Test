package handlers

import (
	"car_rental_test/modules/car/v1/models"
	"context"
)

type ICarSvc interface {
	GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.BookingQueryParams, error)
	CreateBooking(ctx context.Context, input models.Booking) error
	UpdateBooking(ctx context.Context, id int64, input models.Booking) error
	DeleteBooking(ctx context.Context, id int64) (int, error)
}

type Handler struct {
	carSvc ICarSvc
}

func New(carSvc ICarSvc) *Handler {
	return &Handler{
		carSvc: carSvc,
	}
}
