package services

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
)

// Get Bookings by Parameters
func (s *Service) GetBookingsByParams(ctx context.Context, params models.BookingQueryParams) ([]models.BookingQueryParams, error) {
	bookings, err := s.carData.GetBookingsByParams(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetBookingsByParams]")
	}
	return bookings, nil
}

// Create Booking with total cost calculation
func (s *Service) CreateBooking(ctx context.Context, booking models.Booking) error {
	rentPrice, err := s.carData.GetCarRentPrice(ctx, int64(booking.CarsID))
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][GetCarRentPrice]")
	}

	days := int64(booking.EndTime.Sub(booking.StartTime).Hours()/24) + 1
	booking.TotalCost = days * rentPrice

	_, err = s.carData.CreateBooking(ctx, booking)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][CreateBooking][CreateBooking]")
	}

	return nil
}

// Update Booking with total cost recalculation
func (s *Service) UpdateBooking(ctx context.Context, id int64, booking models.Booking) error {
	rentPrice, err := s.carData.GetCarRentPrice(ctx, int64(booking.CarsID))
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateBooking][GetCarRentPrice]")
	}

	days := int(booking.EndTime.Sub(booking.StartTime).Hours()/24) + 1
	booking.TotalCost = int64(days) * rentPrice

	_, err = s.carData.UpdateBooking(ctx, id, booking)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][UpdateBooking][UpdateBooking]")
	}

	return nil
}

// Delete Booking
func (s *Service) DeleteBooking(ctx context.Context, id int64) (int, error) {
	rowsAffected, err := s.carData.DeleteBooking(ctx, id)
	if err != nil {
		return 0, errors.Wrap(err, "[SERVICE][DeleteBooking]")
	}
	return rowsAffected, nil
}
