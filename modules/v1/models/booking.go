package models

type Booking struct {
	BookingID  int    `json:"booking_id"`
	CustomerID int    `json:"customer_id"`
	CarsID     int    `json:"cars_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	TotalCost  int64  `json:"total_cost"`
	Finished   bool   `json:"finished"`
}

type BookingQueryParams struct {
	CustomerID *int    `form:"customer_id"`
	CarsID     *int    `form:"cars_id"`
	StartTime  *string `form:"start_time"`
	EndTime    *string `form:"end_time"`
	Finished   *bool   `form:"finished"`
}
