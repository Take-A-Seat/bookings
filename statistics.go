package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"sort"
)

func rankMapStringInt(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})
	ranked := make([]string, len(values))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}

func getStatistics(restaurantId primitive.ObjectID) (models.StatisticReservations, error) {
	var statistics models.StatisticReservations
	listReservations, err := getAllReservationsByRestaurantId(restaurantId)
	if err != nil {
		return statistics, err
	}

	persons := make(map[string][]int)
	totalPay := make(map[string][]float64)
	totalMoneyReceived := make(map[string]float64)
	numberReservations := make(map[string]int)
	numberPeopleReturned := make(map[string]int)
	declined := make(map[string]int)
	finish := make(map[string]int)

	listData := make(map[string]int)
	index := 0
	var dateString string
	for indexRes, reservation := range listReservations {
		dateString = reservation.StartReservationDate.Format("2006-01-02")
		if _, ok := listData[dateString]; ok {
			persons[dateString] = append(persons[dateString], reservation.Persons)
			totalMoneyReceived[dateString] += reservation.TotalToPay
			if reservation.Status == "Active" || reservation.Status == "Finished" {
				totalPay[dateString] = append(totalPay[dateString], reservation.TotalToPay)
			}
			numberReservations[dateString] += 1
			if checkReturn(listReservations, reservation.Email, indexRes) == true {
				numberPeopleReturned[dateString] += 1
			}
			if reservation.Status == "Declined" {
				declined[dateString] += 1
			} else if reservation.Status == "Finished" {
				finish[dateString] += 1
			}
		} else {
			listData[dateString] = index
			index += 1
			persons[dateString] = append(persons[dateString], reservation.Persons)
			if reservation.Status == "Active" || reservation.Status == "Finished" {
				totalPay[dateString] = append(totalPay[dateString], reservation.TotalToPay)
			}
			numberReservations[dateString] = 1
			numberPeopleReturned[dateString] = 0
			totalMoneyReceived[dateString] = reservation.TotalToPay
			if reservation.Status == "Declined" {
				declined[dateString] = 1
			} else if reservation.Status == "Finished" {
				finish[dateString] = 1
			}
		}
	}

	// make an array of type string to store our keys
	var keys []string

	// iterate over the map and append all keys to our
	// string array of keys
	for key := range listData {
		keys = append(keys, key)
	}

	// use the sort method to sort our keys array
	sort.Strings(keys)
	for _, key := range keys {
		var chartPersons models.ChartData
		var chartTotalPay models.ChartData
		var chartNumberReservations models.CharWithValue
		var chartNumberPeopleReturned models.CharWithValue
		var chartDeclined models.CharWithValue
		var chartFinished models.CharWithValue
		var chartTotalMoneyReviced models.CharWithValue

		chartPersons.Name = key
		chartTotalPay.Name = key
		chartNumberReservations.Name = key
		chartNumberPeopleReturned.Name = key
		chartDeclined.Name = key
		chartFinished.Name = key
		chartTotalMoneyReviced.Name = key
		var sumInt int
		for index, person := range persons[key] {
			if index == 0 {
				chartPersons.Min = float64(person)
			}
			if person < int(chartPersons.Min) {
				chartPersons.Min = float64(person)
			}
			if person > int(chartPersons.Max) {
				chartPersons.Max = float64(person)
			}
			sumInt += person
		}
		chartPersons.Avg = roundTo(float64(sumInt)/float64(len(persons[key])), 2)
		statistics.Persons = append(statistics.Persons, chartPersons)

		var sumFloat float64
		for index, pay := range totalPay[key] {
			if index == 0 {
				chartTotalPay.Min = pay
			}
			if pay < chartTotalPay.Min {
				chartTotalPay.Min = pay
			}
			if pay > chartTotalPay.Max {
				chartTotalPay.Max = pay
			}
			sumFloat += pay
		}

		if len(totalPay[key]) > 0 {
			chartTotalPay.Avg = roundTo(sumFloat/float64(len(totalPay[key])), 2)
		}
		chartDeclined.Value = float64(declined[key])
		chartFinished.Value = float64(finish[key])
		chartNumberPeopleReturned.Value = float64(numberPeopleReturned[key])
		chartFinished.Value = float64(finish[key])
		chartNumberReservations.Value = float64(numberReservations[key])
		chartTotalMoneyReviced.Value = totalMoneyReceived[key]

		statistics.Declined = append(statistics.Declined, chartDeclined)
		statistics.NumberPeopleReturned = append(statistics.NumberPeopleReturned, chartNumberPeopleReturned)
		statistics.Finished = append(statistics.Finished, chartFinished)
		statistics.TotalPay = append(statistics.TotalPay, chartTotalPay)
		statistics.NumberReservations = append(statistics.NumberReservations, chartNumberReservations)
		statistics.TotalMoneyReceived = append(statistics.TotalMoneyReceived, chartTotalMoneyReviced)
	}

	return statistics, nil
}

func roundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

func checkReturn(listReservation []models.Reservation, email string, indexMaxFind int) bool {
	count := 0
	for indexMax, reservation := range listReservation {
		if indexMax > indexMaxFind {
			return false
		} else {
			if reservation.Email == email {
				count += 1
			}
			if count > 1 {
				return true
			}
		}
	}

	return false
}

func getAllReservationsByRestaurantId(restaurantId primitive.ObjectID) ([]models.Reservation, error) {
	var listReservations []models.Reservation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listReservations, err
	}

	bookingsCollection := client.Database(mongoDatabase).Collection("bookings")
	filter := bson.M{"restaurantId": restaurantId}

	countReservation, err := bookingsCollection.CountDocuments(context.Background(), filter)
	if countReservation > 0 {
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"startReservationDate", 1}})

		cursor, err := bookingsCollection.Find(context.Background(), filter, findOptions)
		if err != nil {
			return listReservations, err
		}

		for cursor.Next(context.TODO()) {
			var booking models.Reservation
			err = cursor.Decode(&booking)
			if err != nil {
				return listReservations, err
			}

			listReservations = append(listReservations, booking)
		}
	}
	return listReservations, nil
}
