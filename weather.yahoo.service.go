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
	var (
		deserializeError = fmt.Errorf("Could not deserialize response")
	)
	var (
		request         *http.Request
		err             error
		output          map[string]interface{}
		queryObject     map[string]interface{}
		resultsObject   map[string]interface{}
		channelObject   map[string]interface{}
		windObject      map[string]interface{}
		conditionObject map[string]interface{}
		itemObject      map[string]interface{}
		speedStrValue   string
		tempStrValue    string
		speedValue      float64
		tempValue       float64
		ok              bool
	)

	//Make API Request
	request, err = http.NewRequest("GET", "https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%2C%20wind%20from%20weather.forecast%20 where%20woeid%20%3D%201105779&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys", nil)
	jsonDecoder := json.NewDecoder(request.Body)
	output = make(map[string]interface{})
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

	speedValue, err = strconv.ParseFloat(speedStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}
	tempValue, err = strconv.ParseFloat(tempStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}

	return &WeatherReport{
		WindSpeed:          float32(speedValue),
		TemperatureCelcius: float32(tempValue),
	}, nil
}
