package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClients struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClients {
    return &OpenWeatherClients{
        apiKey: apiKey,
    }
}

func (o OpenWeatherClients) Coordinates(city string) (Coordinate, error) {
    url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"

    resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))
    if err != nil {
        return Coordinate{}, fmt.Errorf("error get Coordinates: %w", err)
    }

    if resp.StatusCode != 200 {
        return Coordinate{}, fmt.Errorf("error fail get Coordinates: %d", resp.StatusCode)
    }

    var coordinatesResponse []CoordinateResponse
    err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
    if err != nil {
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
