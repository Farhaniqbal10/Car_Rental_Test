package data

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"strings"
)

const (
	getBookingsByParams  = "GetBookingsByParams"
	qGetBookingsByParams = `SELECT * FROM booking WHERE customer_id = $1`
	getBookings          = "GetBookings"
	qGetBookings         = `SELECT booking_id, customer_id, cars_id, start_time, end_time, total_cost, finished FROM booking`
	getBookingByID       = "GetBookingByID"
	qGetBookingByID      = `SELECT booking_id, customer_id, cars_id, start_time, end_time, total_cost, finished FROM booking WHERE booking_id = $1`
	insertBooking        = "InsertBooking"
	qInsertBooking       = `INSERT INTO booking (customer_id, cars_id, start_time, end_time, total_cost, finished) VALUES ($1, $2, $3, $4, $5, $6) RETURNING booking_id`
	updateBooking        = "UpdateBooking"
	qUpdateBooking       = `UPDATE booking SET customer_id = $1, cars_id = $2, start_time = $3, end_time = $4, total_cost = $5, finished = $6 WHERE booking_id = $7`
	deleteBooking        = "DeleteBooking"
	qDeleteBooking       = `DELETE FROM booking WHERE booking_id = $1`
	getCarRentPrice      = "GetCarRentPrice"
	qGetCarRentPrice     = `SELECT rent_price_daily FROM cars WHERE cars_id = $1`
)

var (
	BookingStmts = []statement{
		{getBookingsByParams, qGetBookingsByParams},
		{getBookings, qGetBookings},
		{getBookingByID, qGetBookingByID},
		{insertBooking, qInsertBooking},
		{updateBooking, qUpdateBooking},
		{deleteBooking, qDeleteBooking},
		{getCarRentPrice, qGetCarRentPrice},
	}
)

// func (d *Data) GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.BookingQueryParams, error) {
// 	var bookings []models.BookingQueryParams

// 	err := d.db.SelectContext(ctx, &bookings, "SELECT * FROM booking WHERE customer_id = $1", query.CustomerID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bookings, nil
// }

func (d *Data) GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.BookingQueryParams, error) {
	var (
		bookings   []models.BookingQueryParams
		conditions []string
		args       []interface{}
	)

	queryStr := "SELECT * FROM booking"

	if query.BookingID != 0 {
		conditions = append(conditions, "booking_id = $1")
		args = append(args, query.BookingID)
	}
	if query.CustomerID != nil {
		conditions = append(conditions, "customer_id = $1")
		args = append(args, *query.CustomerID)
	}
	if query.CarsID != nil {
		conditions = append(conditions, "cars_id = $1")
		args = append(args, *query.CarsID)
	}
	if query.StartTime != nil {
		conditions = append(conditions, "start_time = $1")
		args = append(args, *query.StartTime)
	}
	if query.EndTime != nil {
		conditions = append(conditions, "end_time = $1")
		args = append(args, *query.EndTime)
	}
	if query.TotalCost != 0 {
		conditions = append(conditions, "total_cost = $1")
		args = append(args, query.TotalCost)
	}
	if query.Finished != nil && query.Finished.Valid {
		conditions = append(conditions, "finished = $1")
		args = append(args, query.Finished.Bool)
	}

	// Gabungkan kondisi dengan operator AND
	if len(conditions) > 0 {
		queryStr += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Eksekusi query dengan parameter
	err := d.db.SelectContext(ctx, &bookings, queryStr, args...)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (d *Data) CreateBooking(ctx context.Context, booking models.Booking) (int, error) {
	var id int
	err := d.stmt[insertBooking].GetContext(ctx, &id,
		booking.CustomerID, booking.CarsID, booking.StartTime, booking.EndTime, booking.TotalCost, booking.Finished)
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][CreateBooking]")
	}
	return id, nil
}

func (d *Data) GetBookingByID(ctx context.Context, id int) (models.Booking, error) {
	var booking models.Booking
	err := d.stmt[getBookingByID].GetContext(ctx, &booking, id)
	if err != nil {
		return booking, errors.Wrap(err, "[DATA][GetBookingByID]")
	}
	return booking, nil
}

func (d *Data) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	var bookings []models.Booking
	err := d.stmt[getBookings].SelectContext(ctx, &bookings)
	if err != nil {
		return bookings, errors.Wrap(err, "[DATA][GetAllBookings]")
	}
	return bookings, nil
}

func (d *Data) UpdateBooking(ctx context.Context, id int64, input models.Booking) (int, error) {
	query := `UPDATE booking SET customer_id = $1, cars_id = $2, start_time = $3, end_time = $4, finished = $5 WHERE booking_id = $6`

	res, err := d.db.ExecContext(ctx, query, input.CustomerID, input.CarsID, input.StartTime, input.EndTime, input.Finished, id)
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][UpdateBooking]")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][UpdateBooking][RowsAffected]")
	}

	return int(rowsAffected), nil
}

func (d *Data) DeleteBooking(ctx context.Context, id int64) (int, error) {
	result, err := d.db.ExecContext(ctx, "DELETE FROM booking WHERE booking_id = $1", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (d *Data) GetCarRentPrice(ctx context.Context, carID int64) (int64, error) {
	var rentPrice int64
	err := d.stmt[getCarRentPrice].GetContext(ctx, &rentPrice, carID)
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetCarRentPrice]")
	}
	return rentPrice, nil
}
