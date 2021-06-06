package main

import (
	"context"
	"fmt"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func getFreeHours(c *gin.Context,restaurantId string, date string, persons string) error {
	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return err
	}
	numberPersons, err := strconv.Atoi(persons)
	if err != nil {
		return err
	}

	timeT, _ := time.Parse("2006-01-02", date)

	listReservation, err := getBookings(restaurantIdObject, timeT)
	if err != nil {
		return err
	}

	restaurant,err := getRestaurantById(c,restaurantId)
	if err != nil {
		return err
	}

	weekday := timeT.Weekday()
	var programDay models.Program
	for _,program := range restaurant.RestaurantDetails.Program{
		if (program.Day)%7 == int(weekday){
			programDay = program
		}
	}

	fmt.Println(timeT, restaurantIdObject, numberPersons, listReservation,restaurant,programDay)

	return nil
}

func createBooking(booking models.Reservation) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	booking.Id = primitive.NewObjectID()
	_, err = bookingsCollection.InsertOne(context.Background(), bson.M{
		"_id":             booking.Id,
		"persons":         booking.Persons,
		"reservationDate": booking.ReservationDate,
		"restaurantId":    booking.RestaurantId,
		"phone":           booking.Phone,
		"firstName":       booking.FirstName,
		"lastName":        booking.LastName,
		"email":           booking.Email,
		"details":         booking.Details,
		"status":          "Pending",
	})

	return nil
}

func getBookings(restaurantId primitive.ObjectID, date time.Time) ([]models.Reservation, error) {
	var listBookings []models.Reservation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listBookings, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	filter := bson.M{"reservationDate": bson.M{
		"$gte": primitive.NewDateTimeFromTime(date),
		"$lte": primitive.NewDateTimeFromTime(date.Add(time.Hour * 24)),
	}}

	countReservation, err := bookingsCollection.CountDocuments(context.Background(), filter)
	if countReservation > 0 {
		cursor, err := bookingsCollection.Find(context.Background(), filter)
		if err != nil {
			return listBookings, err
		}

		for cursor.Next(context.TODO()) {
			var booking models.Reservation
			err = cursor.Decode(&booking)
			if err != nil {
				return listBookings, err
			}

			listBookings = append(listBookings, booking)
		}
	} else {
		return listBookings, nil
	}

	return listBookings, nil
}
