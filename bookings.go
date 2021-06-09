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
	"math/rand"
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

	listReservation, err := getBookings(restaurantIdObject, timeT, "Accepted")
	if err != nil {
		return availableData, err
	}

	startProgram, err := dateparse.ParseAny(date + " " + programDay.StartAt)
	if err != nil {
		return availableData, err
	}

	endProgram, err := dateparse.ParseAny(date + " " + programDay.EndAt)
	if err != nil {
		return availableData, err
	}

	if time.Now().After(endProgram) == true {
		return availableData, errors.New("Closed")
	}

	var currentTimeToAdd time.Time
	if time.Now().Before(startProgram) == true {
		currentTimeToAdd = startProgram
	} else {
		now := time.Now()
		currentTimeToAdd = time.Date(now.Year(), time.Month(int(now.Month())), now.Day(), now.Hour()+3, 0, 0, 0, startProgram.Location())
	}

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

func getReservationByDateAndStatusAndRestaurantId(date time.Time, listReservation []models.Reservation) error {

	return nil
}

func getAllBookingsByRestaurantAndDate(restaurantId string, date string, filterReq string) ([]models.Reservation, error) {
	var listBookings []models.Reservation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listBookings, err
	}

	timeT, err := time.Parse("2006-01-02", date)
	if err != nil {
		return listBookings, err
	}

	restaurantIdObj, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return listBookings, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	var filter bson.M
	if filterReq=="All"{
		filter = bson.M{"startReservationDate": bson.M{
			"$gte": primitive.NewDateTimeFromTime(timeT),
			"$lte": primitive.NewDateTimeFromTime(timeT.Add(time.Hour * 24)),
		}, "restaurantId": restaurantIdObj}
	}else{
		filter = bson.M{"startReservationDate": bson.M{
			"$gte": primitive.NewDateTimeFromTime(timeT),
			"$lte": primitive.NewDateTimeFromTime(timeT.Add(time.Hour * 24)),
		}, "restaurantId": restaurantIdObj, "status": filterReq}
	}


	countReservation, err := bookingsCollection.CountDocuments(context.Background(), filter)
	if countReservation > 0 {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"startReservationDate", -1}, {"status", 1}})

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

func findDateInListReservation(date time.Time, list []models.Reservation) bool {
	//for _, book := range list {
	//	if date.Equal(book.ReservationDate) == true {
	//		return true
	//	}
	//}
	return false
}

func createBooking(booking models.Reservation, restaurant models.RestaurantWithDetails) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	booking.Id = primitive.NewObjectID()
	_, err = bookingsCollection.InsertOne(context.Background(), bson.M{
		"_id":                  booking.Id,
		"persons":              booking.Persons,
		"startReservationDate": booking.StartReservationDate,
		"endReservationDate":   booking.EndReservationDate,
		"restaurantId":         booking.RestaurantId,
		"phone":                booking.Phone,
		"firstName":            booking.FirstName,
		"lastName":             booking.LastName,
		"email":                booking.Email,
		"details":              booking.Details,
		"status":               "Pending",
	})
	if err != nil {
		return err
	}

	startHourString := strconv.Itoa(booking.StartReservationDate.Hour()) + ":" + strconv.Itoa(booking.StartReservationDate.Minute())
	if strconv.Itoa(booking.StartReservationDate.Minute())=="0"{
		startHourString+="0"
	}
	startDate := strconv.Itoa(booking.StartReservationDate.Year()) + "-" + strconv.Itoa(int(time.Month(int(booking.StartReservationDate.Month())))) + "-" + strconv.Itoa(booking.StartReservationDate.Day())

	endHourString := strconv.Itoa(booking.EndReservationDate.Hour()) + ":" + strconv.Itoa(booking.EndReservationDate.Minute())
	if strconv.Itoa(booking.EndReservationDate.Minute()) =="0"{
		endHourString+="0"
	}
	sendConfirmationCreateReservation(booking.Email, booking.FirstName, restaurant.RestaurantDetails.Name, startHourString, startDate,endHourString,)
	return nil
}

func getBookings(restaurantId primitive.ObjectID, date time.Time, status string) ([]models.Reservation, error) {
	var listBookings []models.Reservation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listBookings, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	filter := bson.M{"startReservationDate": bson.M{
		"$gte": primitive.NewDateTimeFromTime(date),
		"$lte": primitive.NewDateTimeFromTime(date.Add(time.Hour * 24)),
	}, "status": status, "restaurantId": restaurantId}

	countReservation, err := bookingsCollection.CountDocuments(context.Background(), filter)
	if countReservation > 0 {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"startReservationDate", 1}})

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func confirmBooking(booking models.Reservation, c *gin.Context, bookingId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	booking.Code = RandStringBytes(6)
	bookingCollection := client.Database(mongoDatabase).Collection("bookings")
	bookingIdObj, err := primitive.ObjectIDFromHex(bookingId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": bookingIdObj}

	updateObject := bson.D{{"$set", bson.D{
		{"status", "Accepted"},
		{"tableId", booking.TableId},
		{"messageToClient", booking.MessageToClient},
		{"code", booking.Code},
	}}}
	_, err = bookingCollection.UpdateOne(context.Background(), filter, updateObject)
	if err != nil {
		return err
	}

	bookingDb, err := getBookingById(bookingIdObj)
	if err != nil {
		return err
	}

	restaurant, err := getRestaurantById(c, booking.RestaurantId.Hex())
	if err != nil {
		return err
	}

	sendConfirmationAcceptReservation(bookingDb.Email, bookingDb.FirstName, booking.MessageToClient, restaurant.RestaurantDetails.Name, booking.Code)
	return nil
}

func getBookingById(bookingId primitive.ObjectID) (models.Reservation, error) {
	var booking models.Reservation
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return booking, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	filter := bson.M{"_id": bookingId}
	err = bookingsCollection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return booking, err

	}
	return booking, nil
}

func availableTables(restaurantId string, startReservation string, endReservation string) error {
	starT, err := time.Parse("2006-01-02T15:04", startReservation)
	if err != nil {
		return err
	}
	endT, err := time.Parse("2006-01-02T15:04", endReservation)
	if err != nil {
		return err
	}

	tableId, err := primitive.ObjectIDFromHex("60bfb0e1cc5c5c17a782aa1c")
	if err != nil {
		return err
	}

	check, err := checkTableAvailableInInterval(tableId, starT, endT)
	if err != nil {
		return err
	}

	fmt.Println(check)

	return nil
}

func checkTableAvailableInInterval(tableId primitive.ObjectID, start time.Time, end time.Time) (bool, error) {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return false, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	//filterObject := bson.M{
	//	"$or": bson.A{bson.M{"associationId": associationObjectId}, bson.M{"associationId": primitive.NilObjectID}},
	//}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"startReservationDate": bson.M{"$gte": primitive.NewDateTimeFromTime(start), "$lt": primitive.NewDateTimeFromTime(end)},
				"endReservationDate": bson.M{"$gte": primitive.NewDateTimeFromTime(end), "$lt": primitive.NewDateTimeFromTime(start)},
			},
			bson.M{"startReservationDate": bson.M{"$lte": primitive.NewDateTimeFromTime(start), "$lt": primitive.NewDateTimeFromTime(end)},
				"endReservationDate": bson.M{"$gt": primitive.NewDateTimeFromTime(end), "$gte": primitive.NewDateTimeFromTime(start)},
			},
			bson.M{"startReservationDate": bson.M{"$gte": primitive.NewDateTimeFromTime(start), "$lte": primitive.NewDateTimeFromTime(end)},
				"endReservationDate": bson.M{"$lt": primitive.NewDateTimeFromTime(end), "$lte": primitive.NewDateTimeFromTime(start)},
			},
		}, "status": "Accepted", "tableId": tableId}

	countReservation, err := bookingsCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	if countReservation > 0 {
		return true, nil
	}

	return false, nil
}
