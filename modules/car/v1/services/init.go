package services

import (
	"car_rental_test/modules/car/v1/models"
	"context"
)

type ICarData interface {
	GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.BookingQueryParams, error)
	CreateBooking(ctx context.Context, input models.Booking) (int, error)
	UpdateBooking(ctx context.Context, id int64, input models.Booking) (int, error)
	DeleteBooking(ctx context.Context, id int64) (int, error)
	GetCarRentPrice(ctx context.Context, carID int64) (int64, error) // Fungsi baru untuk mendapatkan harga sewa mobil
}

type Service struct {
	carData ICarData
}

func New(carData ICarData) *Service {
	return &Service{
		carData: carData,
	}
}
