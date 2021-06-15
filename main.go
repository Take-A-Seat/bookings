package main

import (
	"flag"
	"github.com/Take-A-Seat/auth/validatorAuth"
	"github.com/Take-A-Seat/storage/ws"
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

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	port := os.Getenv("TAKEASEAT_RESERVATION_PORT")
	if port == "" {
		port = "9220"
	}

	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendAllService()
	go ws.WebsocketManager.SendGroupService()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()


	router.Use(cors.New(cors.Config{
		AllowOrigins:           []string{"*"},
		AllowMethods:           []string{"PUT", "PATCH", "DELETE", "GET", "POST", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Accepts", "Connection", "Authorization", "X-Requested-With", "Sec-WebSocket-Protocol", "Sec-WebSocket-Key"},
		ExposeHeaders:          []string{"Content-Length"},
		MaxAge:                 1 * time.Minute,
		AllowWebSockets:        true,
		AllowCredentials:       true,
		AllowBrowserExtensions: true,
	}))

	freeRoute := router.Group("/bookings")
	{
		freeRoute.GET("/ws/:channel", ws.WebsocketManager.WsClient)
		freeRoute.POST("/restaurant/:restaurantId/booking", handlerCreateBooking)
		freeRoute.POST("/restaurant/:restaurantId/booking/:bookingId/confirm", handlerCreateBooking)
		freeRoute.PUT("/id/:bookingId/status", handleUpdateStatusReservation)
		freeRoute.PUT("/id/:bookingId/products", handlerUpdateProductsFromBooking)
		freeRoute.PUT("/id/:bookingId/assistance", handlerUpdateAssistanceFromBooking)
		freeRoute.GET("/restaurant/:restaurantId/email/:email/code/:code", handlerGetBookingByIdUser)
		freeRoute.GET("/restaurant/:restaurantId/booking-hours/date/:date/persons/:persons", handlerGetBookingAvailable)
		freeRoute.GET("/restaurant/:restaurantId/dataInterval/date/:date", handlerGetDataIntervals)
	}

	protectedUsers := router.Group("/bookings")
	protectedUsers.Use(validatorAuth.AuthMiddleware(apiUrl + "/auth/isAuthenticated"))
	{
		protectedUsers.GET("/restaurant/:restaurantId/date/:date/:filter", handlerGetBookingByRestaurantAndDate)
		protectedUsers.GET("/restaurant/:restaurantId/availableTables/:startRes/:endRes", handleGetAvailableTables)
		protectedUsers.GET("/id/:id", handlerGetBookingByIdManager)
		protectedUsers.GET("/restaurant/:restaurantId/statistics", handlerGetStatisticsByRestaurantId)
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Port already in use!")
	}

}
