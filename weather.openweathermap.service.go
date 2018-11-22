package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		response         *http.Response
		err              error
		output           map[string]interface{}
		mainObject       map[string]interface{}
		windObject       map[string]interface{}
		speedMpsValue    float64
		tempKelvinValue  float64
		tempCelciusValue float64
		ok               bool
	)

	//Make API Request
	response, err = http.Get("http://api.openweathermap.org/data/2.5/weather?q=sydney,AU&appid=2326504fb9b100bee21400190e4dbe6d")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Error: %s", response.Status)
	}
	jsonDecoder := json.NewDecoder(response.Body)
	jsonDecoder.Decode(&output)

	//---------------------------------
	//----------GET WIND SPEED---------
	//---------------------------------
	//wind.speed
	windObject, ok = output["wind"].(map[string]interface{})
	if !ok {
		return nil, deserializeError
	}

	speedMpsValue, ok = windObject["speed"].(float64)
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

	tempKelvinValue, ok = mainObject["temp"].(float64)
	if !ok {
		return nil, deserializeError
	}

	//Convert the temperature to celcius as this service provides it in farenheit
	tempCelciusValue = (tempKelvinValue - 273.15)

	return &WeatherReport{
		WindSpeedMetersPerSecond: float32(speedMpsValue),
		TemperatureCelcius:       float32(tempCelciusValue),
	}, nil
}
