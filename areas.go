package main

import (
	"encoding/json"
	"errors"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getAreasByRestaurantId(c *gin.Context, restaurantId string) ([]models.Area, error) {
	var areas []models.Area

	//set request URL
	requestUrl := apiUrl + "/restaurants/id/" + restaurantId + "/areas"
	//create new http Client
	apiClient := &http.Client{}
	//create new request
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return areas, err
	}

	request.Header.Set("Content-Type", "application/json")
	if c.Request.Header == nil {
		return areas, errors.New("Unauthorized")
	}

	request.Header.Set("Authorization", c.Request.Header["Authorization"][0])
	validateUserResponse, getErr := apiClient.Do(request)
	if getErr != nil {
		return areas, err
	}

	if validateUserResponse.Body != nil {
		defer validateUserResponse.Body.Close()
	}

	userData, userErr := ioutil.ReadAll(validateUserResponse.Body)
	if userErr != nil {
		return areas, userErr
	}

	//parse binary to object
	err = json.Unmarshal(userData, &areas)
	if err != nil {
		panic(err)
	}

	return areas, nil
}
