package services

import (
	"car_rental_test/modules/v1/data"
	"car_rental_test/modules/v1/models"
	"context"
	"errors"
)

type BookingService struct {
	Repo *data.BookingRepository
}

func NewBookingService(repo *data.BookingRepository) *BookingService {
	return &BookingService{Repo: repo}
}

// Get All Bookings
func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	return s.Repo.GetAll(ctx)
}

// Get Booking by ID
func (s *BookingService) GetBookingByID(ctx context.Context, bookingID int) (*models.Booking, error) {
	return s.Repo.GetByID(ctx, bookingID)
}

// GetBookingsByParams retrieves bookings based on optional filters
func (s *BookingService) GetBookingsByParams(ctx context.Context, params models.BookingQueryParams) ([]models.Booking, error) {
	var finished *bool
	if params.Finished != nil { // Pastikan tidak nil
		finished = &params.Finished.Bool
	}

	return s.Repo.GetBookingsByParams(ctx, params.CustomerID, params.CarsID, finished)
}

// Create Booking with total cost calculation
func (s *BookingService) CreateBooking(ctx context.Context, booking models.Booking) (int64, error) {
	rentPrice, err := s.Repo.GetCarRentPrice(ctx, int64(booking.CarsID))
	if err != nil {
		return 0, errors.New("failed to get car price: " + err.Error())
	}

	start := booking.StartTime
	end := booking.EndTime
	days := int64(end.Sub(start).Hours()/24) + 1
	booking.TotalCost = days * rentPrice

	_, err = s.Repo.Create(ctx, &booking)
	if err != nil {
		return 0, err
	}

	return booking.TotalCost, nil
}

// Update Booking with total cost recalculation
func (s *BookingService) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	rentPrice, err := s.Repo.GetCarRentPrice(ctx, int64(booking.CarsID))
	if err != nil {
		return errors.New("failed to get car price")
	}

	days := int(booking.EndTime.Sub(booking.StartTime).Hours()/24) + 1
	booking.TotalCost = int64(days) * rentPrice

	return s.Repo.Update(ctx, booking)
}

// Delete Booking
func (s *BookingService) DeleteBooking(ctx context.Context, bookingID int) error {
	return s.Repo.Delete(ctx, bookingID)
}
