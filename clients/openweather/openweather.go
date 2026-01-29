package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
	apiURL string // https://api.openweathermap.org/data/2.5/weather
	geoURL string // http://api.openweathermap.org/geo/1.0/direct
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
		apiURL: "https://api.openweathermap.org/data/2.5/weather",
		geoURL: "http://api.openweathermap.org/geo/1.0/direct",
	}
}

func (o OpenWeatherClient) Coordinates(ctx context.Context, city string) (Coordinate, error) {
	url := fmt.Sprintf("%s?q=%s&limit=5&appid=%s", o.geoURL, city, o.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Coordinate{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return Coordinate{}, fmt.Errorf("error get Coordinates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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

func (o OpenWeatherClient) Weather(ctx context.Context, lat float64, lon float64) (Weather, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", o.apiURL, lat, lon, o.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Weather{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return Weather{}, fmt.Errorf("error get weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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
