package data

import (
	"car_rental_test/modules/v1/models"
	"context"
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// BookingRepository handles database operations for bookings
type BookingRepository struct {
	db *sqlx.DB
}

// NewBookingRepository initializes a new BookingRepository
func NewBookingRepository(db *sqlx.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// SQL Queries
const (
	insertBookingQuery = `INSERT INTO booking (customer_id, cars_id, start_time, end_time, total_cost, finished) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING booking_id`
	updateBookingQuery  = `UPDATE booking SET customer_id = $1, cars_id = $2, start_time = $3, end_time = $4, total_cost = $5, finished = $6 WHERE booking_id = $7`
	deleteBookingQuery  = `DELETE FROM booking WHERE booking_id = $1`
	getBookingByIDQuery = `SELECT booking_id, customer_id, cars_id, start_time, end_time, total_cost, finished FROM booking WHERE booking_id = $1`
	getAllBookingsQuery = `SELECT booking_id, customer_id, cars_id, start_time, end_time, total_cost, finished FROM booking`
)

// Create inserts a new booking into the database
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, insertBookingQuery,
		booking.CustomerID, booking.CarsID, booking.StartTime, booking.EndTime, booking.TotalCost, booking.Finished).Scan(&id)
	if err != nil {
		return 0, errors.New("[DATA][Create] failed to insert booking: " + err.Error())
	}
	return id, nil
}

// GetBookingsByParams retrieves bookings based on optional filters
func (r *BookingRepository) GetBookingsByParams(ctx context.Context, customerID, carID *int, finished *bool) ([]models.Booking, error) {
	var bookings []models.Booking
	query := `SELECT booking_id, customer_id, cars_id, start_time, end_time, total_cost, finished FROM booking WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if customerID != nil {
		query += ` AND customer_id = $` + strconv.Itoa(argIndex)
		args = append(args, *customerID)
		argIndex++
	}
	if carID != nil {
		query += ` AND cars_id = $` + strconv.Itoa(argIndex)
		args = append(args, *carID)
		argIndex++
	}
	if finished != nil {
		query += ` AND finished = $` + strconv.Itoa(argIndex)
		args = append(args, *finished)
		argIndex++
	}

	err := r.db.SelectContext(ctx, &bookings, query, args...)
	if err != nil {
		return nil, errors.New("[DATA][GetBookingsByParams] failed to fetch bookings: " + err.Error())
	}
	return bookings, nil
}

// Update modifies an existing booking
func (r *BookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	_, err := r.db.ExecContext(ctx, updateBookingQuery,
		booking.CustomerID, booking.CarsID, booking.StartTime, booking.EndTime, booking.TotalCost, booking.Finished, booking.BookingID)
	if err != nil {
		return errors.New("[DATA][Update] failed to update booking: " + err.Error())
	}
	return nil
}

// Delete removes a booking from the database
func (r *BookingRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, deleteBookingQuery, id)
	if err != nil {
		return errors.New("[DATA][Delete] failed to delete booking: " + err.Error())
	}
	return nil
}

// GetByID fetches a booking by its ID
func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.GetContext(ctx, &booking, getBookingByIDQuery, id)
	if err != nil {
		return nil, errors.New("[DATA][GetByID] booking not found: " + err.Error())
	}
	return &booking, nil
}

// GetCarRentPrice retrieves the rent price of a car by its ID
func (r *BookingRepository) GetCarRentPrice(ctx context.Context, carID int64) (int64, error) {
	var rentPrice int64
	query := "SELECT rent_price_daily FROM cars WHERE cars_id = $1"
	err := r.db.GetContext(ctx, &rentPrice, query, carID)
	if err != nil {
		return 0, errors.New("[DATA][GetCarRentPrice] failed to get car rent price: " + err.Error())
	}
	return rentPrice, nil
}

// GetAll retrieves all bookings
func (r *BookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.SelectContext(ctx, &bookings, getAllBookingsQuery)
	if err != nil {
		return nil, errors.New("[DATA][GetAll] failed to fetch booking: " + err.Error())
	}
	return bookings, nil
}
