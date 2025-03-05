package models

import (
	"time"

	"github.com/guregu/null/v5"
)

type BookingQuery struct {
	Keyword string    `form:"keyword"`
	Code    string    `form:"code"`
	Period  null.Time `form:"period"`
}

type Booking struct {
	BookingID   int       `json:"booking_id" db:"booking_id"`
	CustomerID  int       `json:"customer_id" db:"customer_id"`
	CarsID      int       `json:"cars_id" db:"cars_id"`
	StartPeriod time.Time `json:"start_period" db:"start_period"`
	EndPeriod   time.Time `json:"end_period" db:"end_period"`
	TotalCost   int64     `json:"total_cost" db:"total_cost"`
	Finished    bool      `json:"finished" db:"finished"`
}

type BookingQueryParams struct {
	BookingID   int        `form:"booking_id" db:"booking_id"`
	CustomerID  null.Int64 `form:"customer_id" db:"customer_id"`
	CarsID      *int       `form:"cars_id" db:"cars_id"`
	StartPeriod *time.Time `form:"start_period" db:"start_period"`
	EndPeriod   *time.Time `form:"end_period" db:"end_period"`
	TotalCost   int64      `form:"total_cost" db:"total_cost"`
	Finished    *bool      `form:"finished" db:"finished"`
}
