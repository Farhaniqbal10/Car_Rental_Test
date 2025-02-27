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

// CarHandler menyimpan referensi ke database
type CarHandler struct {
	DB *config.Database
}

// Constructor untuk CarHandler
func NewCarHandler(db *config.Database) *CarHandler {
	return &CarHandler{DB: db}
}

// GetCars - Mengambil semua mobil dari database
func (h *CarHandler) GetCarsByParams(c *gin.Context) {
	fmt.Println("Query Params:", c.Request.URL.RawQuery)
	var params models.CarQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter query tidak valid", "detail": err.Error()})
		return
	}

	conditions := []string{}
	values := []interface{}{}
	i := 1

	if params.CarsID != nil {
		conditions = append(conditions, fmt.Sprintf("cars_id = $%d", i))
		values = append(values, *params.CarsID)
		i++
	}
	if params.Name != nil {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", i))
		values = append(values, "%"+*params.Name+"%")
		i++
	}
	if params.RentPriceDaily != nil {
		conditions = append(conditions, fmt.Sprintf("rent_price_daily = $%d", i))
		values = append(values, *params.RentPriceDaily)
		i++
	}
	if params.Stock != nil {
		conditions = append(conditions, fmt.Sprintf("stock = $%d", i))
		values = append(values, *params.Stock)
		i++
	}

	query := "SELECT * FROM cars"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println("Query:", query, "Values:", values)

	rows, err := h.DB.Conn.QueryContext(context.Background(), query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data mobil", "detail": err.Error()})
		return
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(&car.CarsID, &car.Name, &car.RentPriceDaily, &car.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data mobil", "detail": err.Error()})
			return
		}
		cars = append(cars, car)
	}

	c.JSON(http.StatusOK, cars)
}

// CreateCar - Menambahkan mobil baru
func (h *CarHandler) CreateCar(c *gin.Context) {
	var car models.Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := "INSERT INTO cars (name, rent_price_daily, stock) VALUES ($1, $2, $3)"
	params := []interface{}{car.Name, car.RentPriceDaily, car.Stock}

	// Debug log
	fmt.Printf("Executing Query: %s\nWith Params: %v\n", query, params)

	_, err := h.DB.Conn.ExecContext(context.Background(), query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create car",
			"query":   query,
			"params":  params,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Car created successfully"})
}

// UpdateCar - Mengupdate data mobil berdasarkan ID
func (h *CarHandler) UpdateCar(c *gin.Context) {
	id := c.Param("id")
	var car models.Car

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := h.DB.Conn.ExecContext(context.Background(), "UPDATE cars SET name=$1, rent_price_daily=$2, stock=$3 WHERE cars_id=$4",
		car.Name, car.RentPriceDaily, car.Stock, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car updated successfully"})
}

// DeleteCar - Menghapus mobil berdasarkan ID
func (h *CarHandler) DeleteCar(c *gin.Context) {
	id := c.Param("id")

	_, err := h.DB.Conn.ExecContext(context.Background(), "DELETE FROM cars WHERE cars_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}
