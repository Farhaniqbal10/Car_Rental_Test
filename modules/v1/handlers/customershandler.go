package handlers

import (
	"car_rental_test/app/config"
	"car_rental_test/modules/v1/models"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CustomerHandler menyimpan referensi ke database
type CustomerHandler struct {
	DB *config.Database
}

// Constructor untuk CustomerHandler
func NewCustomerHandler(db *config.Database) *CustomerHandler {
	return &CustomerHandler{DB: db}
}

// GetCustomersByParams - Mengambil semua pelanggan dari database
func (h *CustomerHandler) GetCustomersByParams(c *gin.Context) {
	fmt.Println("Query Params:", c.Request.URL.RawQuery)
	var params models.CustomersQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter query tidak valid", "detail": err.Error()})
		return
	}

	conditions := []string{}
	values := []interface{}{}
	i := 1

	if params.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", i))
		values = append(values, *params.CustomerID)
		i++
	}
	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", i))
		values = append(values, "%"+*params.Name+"%")
		i++
	}
	if params.NIK != nil {
		conditions = append(conditions, fmt.Sprintf("nik = $%d", i))
		values = append(values, *params.NIK)
		i++
	}
	if params.PhoneNumber != nil {
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", i))
		values = append(values, *params.PhoneNumber)
		i++
	}

	query := "SELECT * FROM customers"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println("Query:", query, "Values:", values)

	rows, err := h.DB.Conn.Query(context.Background(), query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan", "detail": err.Error()})
		return
	}
	defer rows.Close()

	var customers []models.Customers
	for rows.Next() {
		var customer models.Customers
		if err := rows.Scan(&customer.CustomerID, &customer.Name, &customer.NIK, &customer.PhoneNumber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data pelanggan", "detail": err.Error()})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, customers)
}

// CreateCustomer - Menambahkan pelanggan baru
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer models.Customers
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := "INSERT INTO customers (name, nik, phone_number) VALUES ($1, $2, $3)"
	params := []interface{}{customer.Name, customer.NIK, customer.PhoneNumber}

	// Debug log
	fmt.Printf("Executing Query: %s\nWith Params: %v\n", query, params)

	_, err := h.DB.Conn.Exec(context.Background(), query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create customer",
			"query":   query,
			"params":  params,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully"})
}

// UpdateCustomer - Mengupdate data pelanggan berdasarkan ID
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customers

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := h.DB.Conn.Exec(context.Background(), "UPDATE customers SET name=$1, nik=$2, phone_number=$3 WHERE customer_id=$4",
		customer.Name, customer.NIK, customer.PhoneNumber, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

// DeleteCustomer - Menghapus pelanggan berdasarkan ID
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DB.Conn.Exec(context.Background(), "DELETE FROM customers WHERE customer_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
