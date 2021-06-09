package main

import (
	"github.com/Take-A-Seat/auth/validatorAuth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

var mongoHost = "takeaseat.knilq.mongodb.net"
var mongoUser = "admin"
var mongoPass = "p4r0l4"
var mongoDatabase = "TakeASeat"
var apiUrl = "https://api.takeaseat.site"
var hostname = "https://api.takeaseat.site"
var directoryFiles = "/home/takeaseat/manager/web/files/"

func main() {
	port := os.Getenv("TAKEASEAT_RESERVATION_PORT")
	if port == "" {
		port = "9220"
	}

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accepts", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}))

	freeRoute := router.Group("/bookings")
	{
		freeRoute.POST("/restaurant/:restaurantId/booking", handlerCreateBooking)
		freeRoute.POST("/restaurant/:restaurantId/booking/:bookingId/confirm", handlerCreateBooking)
		freeRoute.POST("/restaurant/:restaurantId/email", handlerSendEmail)
		freeRoute.GET("/restaurant/:restaurantId/booking-hours/date/:date/persons/:persons", handlerGetBookingAvailable)
		//freeRoute.POST("/restaurant/:restaurantId/booking/:bookingId/declineFree", handleAcceptReservation)

	}

	protectedUsers := router.Group("/bookings")
	protectedUsers.Use(validatorAuth.AuthMiddleware(apiUrl + "/auth/isAuthenticated"))
	{
		protectedUsers.POST("/id/:bookingId/confirm", handleAcceptReservation)
		protectedUsers.GET("/restaurant/:restaurantId/date/:date/:filter", handlerGetBookingByRestaurantAndDate)
		protectedUsers.GET("/restaurant/:restaurantId/availableTables//:startRes/:endRes", handleGetAvailableTables)
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Port already in use!")
	}

}
