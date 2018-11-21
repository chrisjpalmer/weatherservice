package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WeatherOpenWeatherMapService struct {
}

func NewWeatherOpenWeatherMapService() *WeatherOpenWeatherMapService {
	newWeatherOpenWeatherMapService := new(WeatherOpenWeatherMapService)
	return newWeatherOpenWeatherMapService
}

func (weatherOpenWeatherMapService *WeatherOpenWeatherMapService) getWeather() (*WeatherReport, error) {
	var (
		deserializeError = fmt.Errorf("Could not deserialize response")
	)
	var (
		request            *http.Request
		err                error
		output             map[string]interface{}
		mainObject         map[string]interface{}
		windObject         map[string]interface{}
		speedStrValue      string
		tempStrValue       string
		speedValue         float64
		tempFarenheitValue float64
		tempCelciusValue   float64
		ok                 bool
	)

	//Make API Request
	request, err = http.NewRequest("GET", "http://api.openweathermap.org/data/2.5/weather?q=sydney,AU&appid=2326504fb9b100bee21400190e4dbe6d", nil)
	jsonDecoder := json.NewDecoder(request.Body)
	jsonDecoder.Decode(&output)

	//---------------------------------
	//----------GET WIND SPEED---------
	//---------------------------------
	//wind.speed
	windObject, ok = output["wind"].(map[string]interface{})
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
	//main.temp
	mainObject, ok = output["main"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	tempStrValue, ok = mainObject["temp"].(string)
	if !ok {
		return nil, deserializeError
	}

	speedValue, err = strconv.ParseFloat(speedStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}
	tempFarenheitValue, err = strconv.ParseFloat(tempStrValue, 32)
	if err != nil {
		return nil, deserializeError
	}

	//Convert the temperature to celcius as this service provides it in farenheit
	tempCelciusValue = (tempFarenheitValue - 32.0) * (5.0 / 9.0)

	return &WeatherReport{
		WindSpeed:          float32(speedValue),
		TemperatureCelcius: float32(tempCelciusValue),
	}, nil
}
