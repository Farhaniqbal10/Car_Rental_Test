package models

import (
	"database/sql"
	"time"
)

type Booking struct {
	BookingID  int          `json:"booking_id" db:"booking_id"`
	CustomerID int          `json:"customer_id" db:"customer_id"`
	CarsID     int          `json:"cars_id" db:"cars_id"`
	StartTime  time.Time    `json:"start_time" db:"start_time"`
	EndTime    time.Time    `json:"end_time" db:"end_time"`
	TotalCost  int64        `json:"total_cost" db:"total_cost"`
	Finished   sql.NullBool `json:"finished" db:"finished"`
}

type BookingQueryParams struct {
	BookingID  int           `form:"booking_id" db:"booking_id"`
	CustomerID *int          `form:"customer_id" db:"customer_id"`
	CarsID     *int          `form:"cars_id" db:"cars_id"`
	StartTime  *time.Time    `form:"start_time" db:"start_time"`
	EndTime    *time.Time    `form:"end_time" db:"end_time"`
	TotalCost  int64         `form:"total_cost" db:"total_cost"`
	Finished   *sql.NullBool `form:"finished" db:"finished"`
}
