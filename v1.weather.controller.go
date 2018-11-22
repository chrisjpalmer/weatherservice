package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type V1WeatherController struct {
	weatherService *WeatherService
}

type V1WeatherControllerGetOutput struct {
	WindSpeed          float32 `json:"wind_speed"`
	TemperatureDegrees float32 `json:"temperature_degrees"`
}

func NewV1WeatherController(weatherService *WeatherService) *V1WeatherController {
	newV1WeatherController := new(V1WeatherController)
	newV1WeatherController.weatherService = weatherService
	return newV1WeatherController
}

func (v1WeatherController *V1WeatherController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		report      *WeatherReportForCache
		output      V1WeatherControllerGetOutput
		jsonEncoder *json.Encoder
	)

	//Handle the error should one occur in this defer function
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorString := fmt.Sprintf("500 - \"%s\"", err.Error())
			w.Write([]byte(errorString))
		}
	}()

	//Call the weather service
	report, err = v1WeatherController.weatherService.GetWeather()
	if err != nil {
		return
	}

	//Upon success, set the output
	output = V1WeatherControllerGetOutput{
		WindSpeed:          report.WindSpeedMetersPerSecond,
		TemperatureDegrees: report.TemperatureCelcius,
	}

	//Encode the response as JSON
	jsonEncoder = json.NewEncoder(w)
	err = jsonEncoder.Encode(output)
	if err != nil {
		return //Here for consistency's sake
	}
}
