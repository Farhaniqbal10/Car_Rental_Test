package models

type CarQueryParams struct {
	CarsID         *int    `form:"cars_id"`
	Name           *string `form:"name"`
	RentPriceDaily *int64  `form:"rent_price_daily"`
	Stock          *int    `form:"stock"`
}

type Car struct {
	CarsID         int    `json:"cars_id"`
	Name           string `json:"name"`
	RentPriceDaily int64  `json:"rent_price_daily"`
	Stock          int    `json:"stock"`
}
