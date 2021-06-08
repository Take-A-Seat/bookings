package main

import (
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handlerSendEmail(c *gin.Context) {

	//sendInviteEmail(primitive.ObjectID{}, "calinciucandrei98@gmail.com", "A66554")

	c.JSON(http.StatusOK, gin.H{"message": "Send!"})
}

func handlerGetBookingAvailable(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	persons := c.Param("persons")
	listAvailableHours,err := getFreeHours(c, restaurantId, date, persons)

	if err == nil {
		c.JSON(http.StatusOK, listAvailableHours)
	} else {
		if err.Error()=="Closed"{
			c.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
		}else{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
	}
}

func handlerCreateBooking(c *gin.Context){
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
		c.JSON(http.StatusBadRequest, gin.H{"message": errInsert})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document types were added successfully!"})
}
