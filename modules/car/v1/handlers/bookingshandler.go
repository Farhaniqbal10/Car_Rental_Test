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
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println("[ERROR] Invalid query params:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
		return
	}
	fmt.Println("param2 : ", params)
	bookings, err := h.carSvc.GetBookingsByParams(ctx, params)
	if err != nil {
		log.Println("[ERROR] Failed to fetch bookings with params:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	fmt.Println("handler2")
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// Create Booking
func (h *Handler) CreateBooking(c *gin.Context) {
	ctx := c.Request.Context()
	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.carSvc.CreateBooking(ctx, booking)
	if err != nil {
		log.Println("[ERROR] Failed to create booking:", err)
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
		log.Println("[ERROR] Invalid booking ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.carSvc.UpdateBooking(ctx, id, booking)
	if err != nil {
		log.Println("[ERROR] Failed to update booking:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}

// Delete Booking
func (h *Handler) DeleteBooking(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] Invalid booking ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	_, err = h.carSvc.DeleteBooking(ctx, id)
	if err != nil {
		log.Println("[ERROR] Failed to delete booking:", err)
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
