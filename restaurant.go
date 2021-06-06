package main

import (
	"encoding/json"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getRestaurantById(c *gin.Context,restaurantId string) (models.RestaurantWithDetails, error) {
	var restaurant models.RestaurantWithDetails

	//set request URL
	requestUrl := apiUrl + "/restaurants/id/"+restaurantId
	//create new http Client
	apiClient := &http.Client{}
	//create new request
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return restaurant, err
	}

	request.Header.Set("Content-Type", "application/json")


	validateUserResponse, getErr := apiClient.Do(request)
	if getErr != nil {
		return restaurant, err
	}

	if validateUserResponse.Body != nil {
		defer validateUserResponse.Body.Close()
	}

	userData, userErr := ioutil.ReadAll(validateUserResponse.Body)
	if userErr != nil {
		return restaurant, userErr
	}

	//parse binary to object
	err = json.Unmarshal(userData, &restaurant)
	if err != nil {
		panic(err)
	}

	return restaurant, nil
}
