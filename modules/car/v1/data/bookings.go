package data

import (
	"car_rental_test/modules/car/v1/models"
	"car_rental_test/pkg/errors"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	getBookingsByParams  = "GetBookingsByParams"
	qGetBookingsByParams = `
SELECT * FROM booking 
WHERE 
    (booking_id = $1 OR $1 = 0)
    AND (customer_id = $2 OR $2 = 0)
    AND (cars_id = $3 OR $3 = 0)
    AND (start_period >= $4 OR $4 = '0001-01-01')
    AND (end_period <= $5 OR $5 = '0001-01-01')
    AND (finished = $6)
ORDER BY booking_id ASC;
`

	getBookings       = "GetBookings"
	qGetBookings      = `SELECT booking_id, customer_id, cars_id, start_period, end_period, total_cost, finished FROM booking ORDER BY booking_id ASC`
	getBookingByID    = "GetBookingByID"
	qGetBookingByID   = `SELECT booking_id, customer_id, cars_id, start_period, end_period, total_cost, finished FROM booking WHERE booking_id = $1`
	insertBooking     = "InsertBooking"
	qInsertBooking    = `INSERT INTO booking (customer_id, cars_id, start_period, end_period, total_cost, finished) VALUES ($1, $2, $3, $4, $5, $6) RETURNING booking_id`
	updateBooking     = "UpdateBooking"
	qUpdateBooking    = `UPDATE booking SET customer_id = $1, cars_id = $2, start_period = $3, end_period = $4, total_cost = $5, finished = $6 WHERE booking_id = $7`
	deleteBooking     = "DeleteBooking"
	qDeleteBooking    = `DELETE FROM booking WHERE booking_id = $1`
	getCarRentPrice   = "GetCarRentPrice"
	qGetCarRentPrice  = `SELECT rent_price_daily FROM cars WHERE cars_id = $1`
	decreaseCarStock  = "DecreaseCarStock"
	qDecreaseCarStock = `UPDATE cars SET stock = stock - 1 WHERE cars_id = $1 AND stock > 0`
	increaseCarStock  = "IncreaseCarStock"
	qIncreaseCarStock = `UPDATE cars SET stock = stock + 1 WHERE cars_id = $1`
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
		{decreaseCarStock, qDecreaseCarStock},
		{increaseCarStock, qIncreaseCarStock},
	}
)

// func (d *Data) GetBookingsByParams(ctx context.Context, query models.BookingQueryParams) ([]models.BookingQueryParams, error) {
// 	var bookings []models.BookingQueryParams

//		err := d.db.SelectContext(ctx, &bookings, "SELECT * FROM booking WHERE customer_id = $1", query.CustomerID)
//		if err != nil {
//			return nil, err
//		}
//		return bookings, nil
//	}
func (d *Data) GetBookingsByParams(ctx context.Context, tx *sqlx.Tx, query models.BookingQueryParams) ([]models.Booking, error) {
	booking := []models.Booking{} //ganti jadi model booking
	fmt.Println("data1")
	stmt := d.stmt[getBookingsByParams]
	if tx != nil {
		stmt = tx.Stmtx(stmt)
	}

	params := []interface{}{
		query.BookingID,
		query.CustomerID,
		query.CarsID,
		query.StartPeriod,
		query.EndPeriod,
		query.Finished,
	}

	err := stmt.SelectContext(ctx, &booking, params...)
	if err != nil {
		return booking, errors.Wrap(err, "[DATA][GetBookingsByParams]")
	}
	fmt.Println("data2")
	return booking, nil
}

func (d *Data) CreateBooking(ctx context.Context, tx *sqlx.Tx, input models.Booking) (int64, error) {

	var id int64
	stmt := d.stmt[insertBooking]
	if tx != nil {
		stmt = tx.StmtxContext(ctx, stmt)
	}

	err := stmt.GetContext(ctx, &id,
		input.CustomerID,
		input.CarsID,
		input.StartPeriod,
		input.EndPeriod,
		input.TotalCost,
		input.Finished)

	if err != nil {
		return id, errors.Wrap(err, "[DATA][insertArea]")
	}

	return id, nil
}

func (d *Data) GetBookingByID(ctx context.Context, id int64) (models.Booking, error) {
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

func (d *Data) UpdateBooking(ctx context.Context, id int64, tx *sqlx.Tx, input models.Booking) (models.Booking, error) {
	bookings := models.Booking{}

	stmt := d.stmt[updateBooking]

	if tx != nil {
		stmt = tx.StmtxContext(ctx, stmt)
	}

	_, err := stmt.ExecContext(ctx,
		input.CustomerID,
		input.CarsID,
		input.StartPeriod,
		input.EndPeriod,
		input.TotalCost,
		input.Finished,
		id)
	if err != nil {
		return bookings, errors.Wrap(err, "[DATA][updateArea]")
	}

	return bookings, nil
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

func (d *Data) DecreaseCarStock(ctx context.Context, tx *sqlx.Tx, carID int64) (int64, error) {
	result, err := tx.ExecContext(ctx, qDecreaseCarStock, carID)
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][DecreaseCarStock]")
	}

	// Ambil jumlah baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][DecreaseCarStock][RowsAffected]")
	}

	return rowsAffected, nil
}

func (d *Data) IncreaseCarStock(ctx context.Context, tx *sqlx.Tx, carID int64) (int64, error) {
	result, err := tx.ExecContext(ctx, qIncreaseCarStock, carID) // Hanya kirim carID
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][IncreaseCarStock]")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][IncreaseCarStock][RowsAffected]")
	}

	return rowsAffected, nil
}
