package main

import (
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func handlerGetBookingAvailable(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	persons := c.Param("persons")
	listAvailableHours,err := getFreeHours(c, restaurantId, date, persons)

	if err == nil {
		c.JSON(http.StatusOK, listAvailableHours)
	} else {
		if err.Error() == "Closed" {
			c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
	}
}

func handlerGetBookingByIdManager(c *gin.Context) {
	bookingId := c.Param("id")
	bookingIdObj, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := getBookingById(bookingIdObj)

	if err == nil {
		c.JSON(http.StatusOK, booking)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func handlerGetBookingByIdUser(c *gin.Context) {
	bookingId := c.Param("id")
	code := c.Param("code")
	bookingIdObj, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := getBookingById(bookingIdObj)

	if err == nil && booking.Code == code {
		c.JSON(http.StatusOK, booking)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func handlerGetDataIntervals(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	dataInterval, err := getAllIntervalDayBasedDay(c, restaurantId, date)

	if err == nil {
		c.JSON(http.StatusOK, dataInterval)
	} else {
		if err.Error() == "Closed" {
			c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
	}
}

func handlerGetBookingByRestaurantAndDate(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	filter := c.Param("filter")
	listBookings, err := getAllBookingsByRestaurantAndDate(restaurantId, date, filter)

	if err == nil {
		c.JSON(http.StatusOK, listBookings)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func handleUpdateStatusReservation(c *gin.Context) {
	var booking models.Reservation

	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookingId := c.Param("bookingId")
	err = updateStatusBooking(booking, c, bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Accept reservation success"})
	}
}
func handleGetAvailableTables(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	startRes := c.Param("startRes")
	endRes := c.Param("endRes")
	listAreas, err := availableTables(restaurantId, startRes, endRes, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, listAreas)
	}
}

func handlerCreateBooking(c *gin.Context) {
	var booking models.Reservation

	errBindJson := c.ShouldBindJSON(&booking)
	if errBindJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errBindJson.Error()})
		return
	}

	restaurant,err := getRestaurantById(c,booking.RestaurantId.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	errInsert := createBooking(booking,restaurant)
	if errInsert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInsert})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document types were added successfully!"})
}
