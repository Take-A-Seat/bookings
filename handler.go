package main

import (
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func handlerSendEmail(c *gin.Context) {

	sendInviteEmail(primitive.ObjectID{}, "calinciuc.andrei@yahoo.com", "A66554")

	c.JSON(http.StatusOK, gin.H{"message": "Send!"})
}

func handlerGetBookingAvailable(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	persons := c.Param("persons")
	err := getFreeHours(restaurantId, date, persons)
	if err == nil {
		c.JSON(http.StatusOK, "OK")
	} else {
		c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
	}
}

func handlerCreateBooking(c *gin.Context){
	var booking models.Reservation

	errBindJson := c.ShouldBindJSON(&booking)
	if errBindJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errBindJson.Error()})
		return
	}
	errInsert := createBooking(booking)

	if errInsert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errInsert})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document types were added successfully!"})
}
