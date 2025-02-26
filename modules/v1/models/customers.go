package models

type Customers struct {
	CustomerID  int    `json:"customer_id"`
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}

type CustomersQueryParams struct {
	CustomerID  *int    `form:"customer_id"`
	Name        *string `form:"name"`
	NIK         *string `form:"nik"`
	PhoneNumber *string `form:"phone_number"`
}
