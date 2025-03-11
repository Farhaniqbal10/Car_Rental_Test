package models

import (
	"time"

	"github.com/guregu/null/v5"
)

type BookingQuery struct {
	Keyword string `from:"keyword"`
}

// Model utama yang mencerminkan tabel booking di database
type Booking struct {
	BookingID   int64     `json:"booking_id" db:"booking_id"`
	CustomerID  int64     `json:"customer_id" db:"customer_id"`
	CarsID      int64     `json:"cars_id" db:"cars_id"`
	StartPeriod time.Time `json:"start_period" db:"start_period"`
	EndPeriod   time.Time `json:"end_period" db:"end_period"`
	TotalCost   int64     `json:"total_cost" db:"total_cost"`
	Finished    bool      `json:"finished" db:"finished"`
}

// Struct untuk pencarian booking berdasarkan parameter tertentu
type BookingQueryParams struct { // ga perlu null
	BookingID   int64     `form:"id"`
	CustomerID  int64     `form:"customer_id"`
	CarsID      int64     `form:"cars_id"`
	StartPeriod time.Time `form:"start_period"`
	EndPeriod   time.Time `form:"end_period"`
	TotalCost   int64     `form:"total_cost"`
	Finished    bool      `form:"finished" db:"finished"`
}

// Struct untuk input dan update booking
type BookingInput struct { //rapihkan
	CustomerID  int64     `json:"customer_id"  binding:"required"`
	CarsID      int64     `json:"cars_id"  binding:"required"`
	StartPeriod time.Time `json:"start_period" binding:"required"`
	EndPeriod   time.Time `json:"end_period" binding:"required"`
}

type BookingUpdate struct {
	BookingID   int64      `json:"booking_id" binding:"required" ` // ID wajib untuk update
	CustomerID  null.Int64 `json:"customer_id" `                   // Bisa null jika tidak di-update
	CarsID      null.Int64 `json:"cars_id"`
	StartPeriod null.Time  `json:"start_period" `
	EndPeriod   null.Time  `json:"end_period" `
	TotalCost   null.Int64 `json:"total_cost" `
	Finished    null.Bool  `json:"finished"`
}

// Struct untuk request penghapusan booking
type BookingDelete struct {
	BookingID int64 `json:"id" form:"id" validate:"required"`
}
