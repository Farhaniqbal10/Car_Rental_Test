package handlers

import (
	"car_rental_test/modules/v1/models"
	"car_rental_test/modules/v1/services"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	Service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{Service: service}
}

// Get Bookings by Params or All Bookings
func (h *BookingHandler) GetBookingsByParams(c *gin.Context) {
	var params models.BookingQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println("[ERROR] Invalid query params:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
		return
	}

	log.Printf("Received query params: %+v\n", params) // Debugging log

	bookings, err := h.Service.GetBookingsByParams(context.Background(), params)
	if err != nil {
		log.Println("[ERROR] Failed to fetch bookings with params:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// Create Booking
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		log.Println("[ERROR] Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	bookingID, err := h.Service.CreateBooking(c.Request.Context(), booking)
	if err != nil {
		log.Println("[ERROR] Failed to create booking:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Booking created successfully",
		"booking_id": bookingID,
	})
}

// Update Booking
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	booking.BookingID = id

	err = h.Service.UpdateBooking(context.Background(), &booking)
	if err != nil {
		log.Println("[ERROR] Failed to update booking:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}

// Delete Booking
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("[ERROR] Invalid booking ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	err = h.Service.DeleteBooking(context.Background(), id)
	if err != nil {
		log.Println("[ERROR] Failed to delete booking:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
