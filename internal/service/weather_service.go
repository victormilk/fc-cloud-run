package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type WeatherService struct {
	APIKey string
	URL    string
}

func NewWeatherService(apiKey string, url string) *WeatherService {
	return &WeatherService{
		APIKey: apiKey,
		URL:    url,
	}
}

type WeatherResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type APIWeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

func (s *WeatherService) GetWeather(city string) (*WeatherResponse, error) {
	baseUrl, err := url.Parse(s.URL)
	if err != nil {
		return nil, errors.New("failed to parse the URL")
	}
	params := url.Values{
		"key": {s.APIKey},
		"q":   {city},
	}
	baseUrl.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Get(baseUrl.String())
	if err != nil {
		log.Println("failed to get weather", err)
		return nil, errors.New("failed to get weather")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get weather")
	}

	var apiResp APIWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, errors.New("an error occurred while decoding the response")
	}

	return &WeatherResponse{
		Celsius:    apiResp.Current.TempC,
		Fahrenheit: apiResp.Current.TempF,
		Kelvin:     celsiusToKelvin(apiResp.Current.TempC),
	}, nil
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
