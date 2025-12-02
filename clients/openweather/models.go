package openweather

type CoordinateResponse struct {
    Name string `json:"name"`
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type Coordinate struct {
    Lat float64
    Lon float64
}
