package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//func getRestaurantById(c *gin.Context,restaurantId string) (models.RestaurantWithDetails, error) {
//	var restaurant models.RestaurantWithDetails
//
//	//set request URL
//	requestUrl := apiUrl + "/restaurants/id/"+restaurantId
//	//create new http Client
//	apiClient := &http.Client{}
//	//create new request
//	request, err := http.NewRequest("GET", requestUrl, nil)
//	if err != nil {
//		return restaurant, err
//	}
//
//	request.Header.Set("Content-Type", "application/json")
//
//
//	validateUserResponse, getErr := apiClient.Do(request)
//	if getErr != nil {
//		return restaurant, err
//	}
//
//	if validateUserResponse.Body != nil {
//		defer validateUserResponse.Body.Close()
//	}
//
//	userData, userErr := ioutil.ReadAll(validateUserResponse.Body)
//	if userErr != nil {
//		return restaurant, userErr
//	}
//
//	//parse binary to object
//	err = json.Unmarshal(userData, &restaurant)
//	if err != nil {
//		panic(err)
//	}
//
//	return restaurant, nil
//}

func getRestaurantById(restaurantId primitive.ObjectID) (models.RestaurantWithDetails, error) {
	var restaurant models.RestaurantWithDetails
	var filter = bson.M{"_id": restaurantId}

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return models.RestaurantWithDetails{}, err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")
	err = restaurantsCollection.FindOne(context.Background(), filter).Decode(&restaurant.RestaurantDetails)

	if err != nil {
		return models.RestaurantWithDetails{}, err
	}

	return restaurant, nil
}
