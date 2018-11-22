package main

import (
	"fmt"
	"time"
)

//Raw weather report object that providers interact with
type WeatherReport struct {
	WindSpeedMetersPerSecond float32
	TemperatureCelcius       float32
}

//For cache storage in the report service
type WeatherReportForCache struct {
	WindSpeedMetersPerSecond float32
	TemperatureCelcius       float32
	LastDownloaded           time.Time
}

type WeatherService struct {
	sydneyCachedWeatherReport    *WeatherReportForCache
	weatherYahooService          *WeatherYahooService
	weatherOpenWeatherMapService *WeatherOpenWeatherMapService
}

func NewWeatherService(weatherYahooService *WeatherYahooService, weatherOpenWeatherMapService *WeatherOpenWeatherMapService) *WeatherService {
	newWeatherService := new(WeatherService)
	newWeatherService.weatherYahooService = weatherYahooService
	newWeatherService.weatherOpenWeatherMapService = weatherOpenWeatherMapService
	return newWeatherService
}

func (weatherService *WeatherService) getRawWeatherFromProvider() (*WeatherReport, error) {
	//First try the yahoo endpoint...
	var (
		weatherReport *WeatherReport
		err           error
	)
	weatherReport, err = weatherService.weatherYahooService.getWeather()
	if err == nil {
		return weatherReport, nil
	}

	weatherReport, err = weatherService.weatherOpenWeatherMapService.getWeather()
	if err == nil {
		return weatherReport, nil
	}

	return nil, fmt.Errorf("Unable to get weather data as neither provider responded: %s", err.Error())
}

func (weatherService *WeatherService) getWeatherFromProvider() (*WeatherReportForCache, error) {
	report, err := weatherService.getRawWeatherFromProvider()
	if err != nil {
		return nil, err
	}
	return &WeatherReportForCache{
		LastDownloaded:           time.Now(),
		WindSpeedMetersPerSecond: report.WindSpeedMetersPerSecond,
		TemperatureCelcius:       report.TemperatureCelcius,
	}, nil
}

func (weatherService *WeatherService) getWeather() (*WeatherReportForCache, error) {
	//Do we have a cached version of the weather report that is younger than 3 seconds?
	if weatherService.sydneyCachedWeatherReport != nil {
		if time.Since(weatherService.sydneyCachedWeatherReport.LastDownloaded).Seconds() <= 3 {
			return weatherService.sydneyCachedWeatherReport, nil
		}
	}

	//We do not have a cached version younger than 3 secs... lets get one
	report, err := weatherService.getWeatherFromProvider()
	if err == nil {
		//Cache it for next time
		weatherService.sydneyCachedWeatherReport = report
		return report, nil
	}

	//We couldn't successfully get the weather... look again in the cache but be less picky this time
	if weatherService.sydneyCachedWeatherReport != nil {
		return weatherService.sydneyCachedWeatherReport, nil
	}

	//We couldn't get anything out of the cache at all... return the error
	return nil, err
}

func (weatherService *WeatherService) GetWeather() (*WeatherReportForCache, error) {
	//Get the weather report
	report, err := weatherService.getWeather()
	if err != nil {
		return nil, err
	}
	return report, nil
}
