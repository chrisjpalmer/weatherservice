package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WeatherYahooService struct {
}

func NewWeatherYahooService() *WeatherYahooService {
	newWeatherYahooService := new(WeatherYahooService)
	return newWeatherYahooService
}

func (weatherYahooService *WeatherYahooService) getWeather() (*WeatherReport, error) {
	//FOR TESTING PURPOSES:
	//return nil, fmt.Errorf("Test break of yahoo service")
	var (
		deserializeError = fmt.Errorf("Could not deserialize response")
	)
	var (
		response            *http.Response
		err                 error
		output              map[string]interface{}
		queryObject         map[string]interface{}
		resultsObject       map[string]interface{}
		channelObject       map[string]interface{}
		windObject          map[string]interface{}
		conditionObject     map[string]interface{}
		itemObject          map[string]interface{}
		speedStrValue       string
		tempStrValue        string
		speedMphValue       float64 //miles per hour
		speedMpsValue       float64 //meters per second
		tempFahrenheitValue float64
		tempCelciusValue    float64
		ok                  bool
	)

	//Make API Request
	response, err = http.Get("https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%2C%20wind%20from%20weather.forecast%20where%20woeid%20%3D%201105779&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %s", response.Status)
	}
	jsonDecoder := json.NewDecoder(response.Body)
	jsonDecoder.Decode(&output)

	//Deserialize the response
	queryObject, ok = output["query"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	resultsObject, ok = queryObject["results"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	channelObject, ok = resultsObject["channel"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	//---------------------------------
	//----------GET WIND SPEED---------
	//---------------------------------
	//query.results.channel.wind.speed
	windObject, ok = channelObject["wind"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	speedStrValue, ok = windObject["speed"].(string)
	if !ok {
		return nil, deserializeError
	}

	//---------------------------------
	//----------GET TEMPERATURE--------
	//---------------------------------
	//query.results.channel.item.condition.temp
	itemObject, ok = channelObject["item"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	conditionObject, ok = itemObject["condition"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	tempStrValue, ok = conditionObject["temp"].(string)
	if !ok {
		return nil, deserializeError
	}

	speedMphValue, err = strconv.ParseFloat(speedStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}

	speedMpsValue = speedMphValue * 0.44704

	tempFahrenheitValue, err = strconv.ParseFloat(tempStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}

	tempCelciusValue = (tempFahrenheitValue - 32.0) * (5.0 / 9.0)

	return &WeatherReport{
		WindSpeedMetersPerSecond: float32(speedMpsValue),
		TemperatureCelcius:       float32(tempCelciusValue),
	}, nil
}
