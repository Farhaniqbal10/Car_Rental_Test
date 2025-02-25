package models

type Customers struct {
	CustomerID  int    `json:"customer_id"`
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}
