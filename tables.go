package main

import (
	"encoding/json"
	"errors"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getTablesByAreaIdAndRestaurantId(c *gin.Context, restaurantId string,areaId string) ([]models.Table, error) {
	var tables []models.Table

	//set request URL
	requestUrl := apiUrl + "/restaurants/id/" + restaurantId + "/areas/"+areaId+"/tables"
	//create new http Client
	apiClient := &http.Client{}
	//create new request
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return tables, err
	}

	request.Header.Set("Content-Type", "application/json")
	if c.Request.Header == nil {
		return tables, errors.New("Unauthorized")
	}

	request.Header.Set("Authorization", c.Request.Header["Authorization"][0])
	validateUserResponse, getErr := apiClient.Do(request)
	if getErr != nil {
		return tables, err
	}

	if validateUserResponse.Body != nil {
		defer validateUserResponse.Body.Close()
	}

	userData, userErr := ioutil.ReadAll(validateUserResponse.Body)
	if userErr != nil {
		return tables, userErr
	}

	//parse binary to object
	err = json.Unmarshal(userData, &tables)
	if err != nil {
		panic(err)
	}

	return tables, nil
}
