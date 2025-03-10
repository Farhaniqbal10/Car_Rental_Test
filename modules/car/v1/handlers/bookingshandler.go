package handlers

import (
	"car_rental_test/modules/car/v1/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get Bookings by Params or All Bookings
func (h *Handler) GetBookingsByParams(c *gin.Context) {
	ctx := c.Request.Context()
	var params models.BookingQueryParams

	fmt.Println("handler1")
	fmt.Println("param1 : ", params)

	// Bind query params ke struct
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println("[ERROR] Invalid query params:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("param2 : ", params)

	// Ambil data bookings berdasarkan params
	bookings, err := h.carSvc.GetBookingsByParams(ctx, params)
	if err != nil {
		log.Println("[ERROR] Failed to fetch bookings with params:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	fmt.Println("handler2")
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// Create Booking
func (h *Handler) CreateBooking(c *gin.Context) {
	ctx := c.Request.Context()

	body := models.BookingInput{}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("[ERROR] Invalid request body:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// if booking.CustomerID == 0 || booking.CarsID == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "CustomerID dan CarsID wajib diisi"})
	// 	return
	// }

	// if booking.StartPeriod.IsZero() {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "StartPeriod wajib diisi"})
	// 	return
	// }

	// if booking.EndPeriod.IsZero() {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "StartPeriod wajib diisi"})
	// 	return
	// }

	booking := models.Booking{
		CustomerID:  body.CustomerID,
		CarsID:      body.CarsID,
		StartPeriod: body.StartPeriod,
		EndPeriod:   body.EndPeriod,
	}

	err := h.carSvc.CreateBooking(ctx, booking)
	if err != nil {
		log.Println("[ERROR] Failed to create booking:", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully"})
}

// Update Booking
func (h *Handler) UpdateBooking(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] Invalid booking ID:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	body := models.BookingUpdate{}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("[ERROR] Invalid request body:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if body.BookingID != id {
		log.Println("[ERROR] Booking ID mismatch:", id, "!=", body.BookingID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID in URL does not match with body"})
		return
	}

	booking := models.Booking{
		BookingID:   body.BookingID,
		CustomerID:  body.CustomerID.ValueOrZero(), // Ambil nilai atau default (0)
		CarsID:      body.CarsID.ValueOrZero(),
		StartPeriod: body.StartPeriod.ValueOrZero(),
		EndPeriod:   body.EndPeriod.ValueOrZero(),
		TotalCost:   body.TotalCost.ValueOrZero(),
		Finished:    body.Finished.Valid && body.Finished.Bool, // Jika Valid, ambil nilai Bool
	}

	err = h.carSvc.UpdateBooking(ctx, id, booking)
	if err != nil {
		log.Println("[ERROR] Failed to update booking:", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}

// Delete Booking
func (h *Handler) DeleteBooking(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] Invalid booking ID:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.carSvc.DeleteBooking(ctx, id)
	if err != nil {
		log.Println("[ERROR] Failed to delete booking:", err.Error())
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
