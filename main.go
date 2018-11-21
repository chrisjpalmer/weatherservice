package main

import (
	"net/http"
)

func main() {
	//Service Dependency Management
	weatherYahooService := NewWeatherYahooService()
	weatherOpenWeatherMapService := NewWeatherOpenWeatherMapService()
	weatherService := NewWeatherService(weatherYahooService, weatherOpenWeatherMapService)

	//Controllers Dependency Management
	v1WeatherController := NewV1WeatherController(weatherService)

	//Install Http Controllers
	http.Handle("/v1/weather", v1WeatherController)
	http.ListenAndServe(":8080", nil)
}
