package handlers

import (
	"car_rental_test/modules/car/v1/models"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCar(c *gin.Context) {
	ctx := c.Request.Context()
	var params models.CarQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println("[ERROR] Invalid query params:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cars, err := h.carSvc.GetCar(ctx, params)
	if err != nil {
		log.Println("[ERROR] Failed to fetch bookings with params:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cars": cars})

}

func (h *Handler) CreateCar(c *gin.Context) {
	ctx := c.Request.Context()

	body := models.CarInput{}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("[ERROR] Invalid request body:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	car := models.Car{
		Name:           body.Name,
		RentPriceDaily: body.RentPriceDaily,
		Stock:          body.Stock,
	}

	_, err := h.carSvc.CreateCar(ctx, car)
	if err != nil {
		log.Println("[ERROR] Failed to create car:", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Car created successfully"})
}

func (h *Handler) UpdateCar(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] Invalid car ID:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	body := models.CarUpdate{}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println("[ERROR] Invalid request body:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if body.CarsID != id {
		log.Println("[ERROR] Cars ID mismatch:", id, "!=", body.CarsID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cars ID in URL does not match with body"})
		return
	}

	car := models.Car{
		CarsID:         body.CarsID,
		Name:           body.Name.ValueOrZero(),
		RentPriceDaily: body.RentPriceDaily.ValueOrZero(),
		Stock:          body.Stock.ValueOrZero(),
	}

	err = h.carSvc.UpdateCar(ctx, id, car)
	if err != nil {
		log.Println("[ERROR] Failed to update Cars:", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cars updated successfully"})
}

func (h *Handler) DeleteCar(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] Invalid Car ID:", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.carSvc.DeleteCar(ctx, id)
	if err != nil {
		log.Println("[ERROR] Failed to delete car:", err.Error())
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}
