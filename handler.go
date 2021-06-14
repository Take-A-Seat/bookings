package main

import (
	"github.com/Take-A-Seat/storage/models"
	"github.com/Take-A-Seat/storage/ws"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)


func handlerGetBookingAvailable(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	date := c.Param("date")
	persons := c.Param("persons")
	listAvailableHours, err := getFreeHours(c, restaurantId, date, persons)

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

func handlerGetStatisticsByRestaurantId(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	restaurantIdObj, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statistics, err := getStatistics(restaurantIdObj)

	if err == nil {
		c.JSON(http.StatusOK, statistics)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func handlerGetBookingByIdUser(c *gin.Context) {
	code := c.Param("code")
	email := c.Param("email")
	restaurantId := c.Param("restaurantId")

	booking, err := getBookingByCodeAndEmail(email, code, restaurantId)
	if err == nil {
		c.JSON(http.StatusOK, booking)
		return
	} else {
		if err.Error() == "Invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		ws.WebsocketManager.SendGroup(booking.RestaurantId.Hex(), ws.RespClient{
			Type: "get_reservations_list",
			Id:   booking.RestaurantId.Hex(),
			Time: time.Now(),
		})

		ws.WebsocketManager.SendGroup(booking.RestaurantId.Hex(), ws.RespClient{
			Type: "get_available_tables",
			Id:   booking.RestaurantId.Hex(),
			Time: time.Now(),
		})

		c.JSON(http.StatusCreated, gin.H{"error": "Update status success"})
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

	restaurant, err := getRestaurantById(c, booking.RestaurantId.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	errInsert := createBooking(booking, restaurant)
	if errInsert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInsert})
		return
	}

	ws.WebsocketManager.SendGroup(booking.RestaurantId.Hex(), ws.RespClient{
		Type: "get_reservations_list",
		Id:   booking.RestaurantId.Hex(),
		Time: time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Document types were added successfully!"})
}

func handlerUpdateProductsFromBooking(c *gin.Context) {
	var booking models.Reservation

	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookingId := c.Param("bookingId")
	err = updateProductsByBookingId(booking, bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ws.WebsocketManager.SendGroup(booking.Id.Hex(), ws.RespClient{
			Type: "update_booking",
			Id:   booking.Id.Hex(),
			Time: time.Now(),
		})

		ws.WebsocketManager.SendGroup(booking.RestaurantId.Hex(), ws.RespClient{
			Type: "get_reservations_list",
			Id:   booking.RestaurantId.Hex(),
			Time: time.Now(),
		})
		c.JSON(http.StatusCreated, gin.H{"error": "Update products from booking success"})
	}
}

func handlerUpdateAssistanceFromBooking(c *gin.Context) {
	var booking models.Reservation

	err := c.ShouldBindJSON(&booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bookingId := c.Param("bookingId")
	err = updateNeedAssistanceByBookingId(booking, bookingId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ws.WebsocketManager.SendGroup(booking.Id.Hex(), ws.RespClient{
			Type: "update_booking",
			Id:   booking.Id.Hex(),
			Time: time.Now(),
		})

		ws.WebsocketManager.SendGroup(booking.RestaurantId.Hex(), ws.RespClient{
			Type: "get_reservations_list",
			Id:   booking.RestaurantId.Hex(),
			Time: time.Now(),
		})
		c.JSON(http.StatusCreated, gin.H{"error": "Update products from booking success"})
	}
}
