package openweather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinate, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"

	resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))
	if err != nil {
		log.Println(err)
		return Coordinate{}, fmt.Errorf("error get Coordinates: %w", err)
	}

	if resp.StatusCode != 200 {
		log.Println(err)
		return Coordinate{}, fmt.Errorf("error fail get Coordinates: %d", resp.StatusCode)
	}

	var coordinatesResponse []CoordinateResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		log.Println(err)
		return Coordinate{}, fmt.Errorf("error unmarshal response: %w", err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinate{}, fmt.Errorf("error emptry Coordinates")
	}

	return Coordinate{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat float64, lon float64) (Weather, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apiKey))
	if err != nil {
		log.Println(err)
		return Weather{}, fmt.Errorf("error get weather: %w", err)
	}

	if resp.StatusCode != 200 {
		log.Println(err)
		return Weather{}, fmt.Errorf("error fail get weather: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		log.Println(err)
		return Weather{}, fmt.Errorf("error unmarshal weather respose: %w", err)
	}

	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil
}
