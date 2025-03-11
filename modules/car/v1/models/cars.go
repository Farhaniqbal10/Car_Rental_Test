package models

import "github.com/guregu/null/v5"

type CarQueryParams struct {
	Keyword        string `form:"keyword"`
	CarsID         int64  `form:"cars_id" json:"cars_id"`
	Name           string `form:"name" json:"name"`
	RentPriceDaily int64  `form:"rent_price_daily" json:"rent_price_daily"`
	Stock          int64  `form:"stock" json:"stock"`
}

type Car struct {
	CarsID         int64  `json:"cars_id" db:"cars_id"`
	Name           string `json:"name" db:"name"`
	RentPriceDaily int64  `json:"rent_price_daily" db:"rent_price_daily"`
	Stock          int64  `json:"stock" db:"stock"`
}

type CarInput struct {
	CarsID         int64  `json:"cars_id" `
	Name           string `json:"name" binding:"required"`
	RentPriceDaily int64  `json:"rent_price_daily" binding:"required"`
	Stock          int64  `json:"stock" binding:"required"`
}

type CarUpdate struct {
	CarsID         int64       `json:"cars_id" binding:"required"`
	Name           null.String `json:"name"`
	RentPriceDaily null.Int64  `json:"rent_price_daily"`
	Stock          null.Int64  `json:"stock"`
}

type CarDelete struct {
	CarsID int64 `json:"id" form:"id" validate:"required"`
}
