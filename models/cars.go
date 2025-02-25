package models

type Car struct {
	CarsID         int    `json:"cars_id"`
	Name           string `json:"name"`
	RentPriceDaily int64  `json:"rent_price_daily"`
	Stock          int    `json:"stock"`
}
