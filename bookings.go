package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

func getFreeHours(c *gin.Context, restaurantId string, date string, persons string) ([]models.AvailableDataReservation, error) {
	var availableData []models.AvailableDataReservation
	var programDay models.Program

	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return availableData, err
	}

	timeT, err := time.Parse("2006-01-02", date)
	if err != nil {
		return availableData, err
	}

	restaurant, err := getRestaurantById(c, restaurantId)
	if err != nil {
		return availableData, err
	}

	weekday := timeT.Weekday()
	for _, program := range restaurant.RestaurantDetails.Program {
		if (program.Day)%7 == int(weekday) {
			programDay = program
		}
	}

	if programDay.Close == true {
		return availableData, errors.New("Closed")

	}

	numberPersons, err := strconv.Atoi(persons)
	if err != nil {
		return availableData, err
	}

	listReservation, err := getBookings(restaurantIdObject, timeT)
	if err != nil {
		return availableData, err
	}

	//layout := "2006-01-02 3:4:5 PM"

	//startProgram, err := time.Parse(layout, date+" "+programDay.StartAt)
	startProgram, err := dateparse.ParseAny(date + " " + programDay.StartAt)
	if err != nil {
		return availableData, err
	}

	endProgram, err := dateparse.ParseAny(date + " " + programDay.EndAt)
	if err != nil {
		return availableData, err
	}

	currentTimeToAdd := startProgram
	var freeDate models.AvailableDataReservation

	for currentTimeToAdd.After(endProgram) == false {
		freeDate.DateTime = currentTimeToAdd
		hour, minutes, _ := currentTimeToAdd.Clock()
		stringHour := strconv.Itoa(hour)
		stringMinutes := strconv.Itoa(minutes)
		if stringMinutes != "0" {
			freeDate.TimeString = stringHour + ":" + stringMinutes
		} else {
			freeDate.TimeString = stringHour + ":" + stringMinutes + "0"

		}
		availableData = append(availableData, freeDate)
		currentTimeToAdd = currentTimeToAdd.Add(time.Minute * 15)
	}

	fmt.Println(timeT, restaurantIdObject, numberPersons, listReservation, restaurant, programDay, startProgram, endProgram)

	return availableData, err
}



func createBooking(booking models.Reservation,restaurant models.RestaurantWithDetails) error {
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
	if err!=nil{
		return err
	}

	sendConfirmationCreateReservation(booking.Email,booking.FirstName,restaurant.RestaurantDetails.Name)
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
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"reservationDate", 1}})

		cursor, err := bookingsCollection.Find(context.Background(), filter, findOptions)
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
