package handlers

import (
	"car_rental_test/app/config"
	"car_rental_test/modules/v1/models"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	DB *config.Database
}

func NewBookingHandler(db *config.Database) *BookingHandler {
	return &BookingHandler{DB: db}
}

func (h *BookingHandler) GetBookingsByParams(c *gin.Context) {
	fmt.Println("Query Params:", c.Request.URL.RawQuery)
	var params models.BookingQueryParams

	// Binding query params ke struct
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
	if params.CarsID != nil {
		conditions = append(conditions, fmt.Sprintf("cars_id = $%d", i))
		values = append(values, *params.CarsID)
		i++
	}
	if params.StartTime != nil {
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", i))
		values = append(values, *params.StartTime)
		i++
	}
	if params.EndTime != nil {
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", i))
		values = append(values, *params.EndTime)
		i++
	}
	if params.Finished != nil {
		conditions = append(conditions, fmt.Sprintf("finished = $%d", i))
		values = append(values, *params.Finished)
		i++
	}

	query := "SELECT * FROM bookings"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println("Query:", query, "Values:", values)

	rows, err := h.DB.Conn.Query(context.Background(), query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data booking", "detail": err.Error()})
		return
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.BookingID, &booking.CustomerID, &booking.CarsID, &booking.StartTime, &booking.EndTime, &booking.TotalCost, &booking.Finished); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data booking", "detail": err.Error()})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}

// Create Booking
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get car rental price
	var rentPrice int64
	err := h.DB.Conn.QueryRow(context.Background(), "SELECT rent_price_daily FROM cars WHERE cars_id=$1", booking.CarsID).Scan(&rentPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get car price", "details": err.Error()})
		return
	}

	// Calculate total cost
	start, _ := time.Parse("2006-01-02", booking.StartTime)
	end, _ := time.Parse("2006-01-02", booking.EndTime)
	days := int64(end.Sub(start).Hours()/24) + 1
	booking.TotalCost = days * rentPrice

	// Insert booking
	query := "INSERT INTO booking (customer_id, cars_id, start_time, end_time, total_cost, finished) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = h.DB.Conn.Exec(context.Background(), query, booking.CustomerID, booking.CarsID, booking.StartTime, booking.EndTime, booking.TotalCost, booking.Finished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully", "total_cost": booking.TotalCost})
}

// Get Bookings
func (h *BookingHandler) GetBookings(c *gin.Context) {
	rows, err := h.DB.Conn.Query(context.Background(), "SELECT * FROM bookings")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.BookingID, &booking.CustomerID, &booking.CarsID, &booking.StartTime, &booking.EndTime, &booking.TotalCost, &booking.Finished); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan bookings"})
			return
		}
		bookings = append(bookings, booking)
	}

	c.JSON(http.StatusOK, bookings)
}

// Update Booking
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Recalculate total cost
	var rentPrice int64
	err := h.DB.Conn.QueryRow(context.Background(), "SELECT rent_price_daily FROM cars WHERE cars_id=$1", booking.CarsID).Scan(&rentPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get car price", "details": err.Error()})
		return
	}

	start, _ := time.Parse("2006-01-02", booking.StartTime)
	end, _ := time.Parse("2006-01-02", booking.EndTime)
	days := int64(end.Sub(start).Hours()/24) + 1
	booking.TotalCost = days * rentPrice

	_, err = h.DB.Conn.Exec(context.Background(), "UPDATE bookings SET customer_id=$1, cars_id=$2, start_time=$3, end_time=$4, total_cost=$5, finished=$6 WHERE booking_id=$7", booking.CustomerID, booking.CarsID, booking.StartTime, booking.EndTime, booking.TotalCost, booking.Finished, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully", "total_cost": booking.TotalCost})
}

// Delete Booking
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id := c.Param("id")
	_, err := h.DB.Conn.Exec(context.Background(), "DELETE FROM bookings WHERE booking_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
