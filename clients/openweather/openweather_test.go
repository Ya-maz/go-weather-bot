package openweather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenWeatherClient_Coordinates(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/geo/1.0/direct", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "Moscow" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"lat": 55.7558, "lon": 37.6173}]`))
		} else if q == "NotFound" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := New("dummy_key")
	client.geoURL = server.URL + "/geo/1.0/direct"

	tests := []struct {
		name    string
		city    string
		wantErr bool
		wantLat float64
	}{
		{"Valid city", "Moscow", false, 55.7558},
		{"Not found", "NotFound", true, 0},
		{"Server error", "Error", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := client.Coordinates(context.Background(), tt.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && coord.Lat != tt.wantLat {
				t.Errorf("got Lat %v, want %v", coord.Lat, tt.wantLat)
			}
		})
	}
}

func TestOpenWeatherClient_Weather(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		if lat == "55.755800" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"main": {"temp": 12.5}}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := New("dummy_key")
	client.apiURL = server.URL + "/data/2.5/weather"

	tests := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
		wantTmp float64
	}{
		{"Valid coords", 55.7558, 37.6173, false, 12.5},
		{"Invalid coords", 0, 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather, err := client.Weather(context.Background(), tt.lat, tt.lon)
			if (err != nil) != tt.wantErr {
				t.Errorf("Weather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && weather.Temp != tt.wantTmp {
				t.Errorf("got Temp %v, want %v", weather.Temp, tt.wantTmp)
			}
		})
	}
}
