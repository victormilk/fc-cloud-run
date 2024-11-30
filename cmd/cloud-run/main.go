package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/victormilk/fc-cloud-run/internal/service"
)

func main() {
	mux := http.NewServeMux()

	addressService := service.NewAddressService(os.Getenv("ADDRESS_API_URL"))
	weatherService := service.NewWeatherService(os.Getenv("WEATHER_API_KEY"), os.Getenv("WEATHER_API_URL"))

	mux.HandleFunc("GET /temp", func(w http.ResponseWriter, r *http.Request) {
		zipCode := r.URL.Query().Get("zipCode")

		address, err := addressService.GetAddress(zipCode)
		if err != nil {
			if errors.Is(err, service.ErrInvalidZipCode) {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, service.ErrNotFoundZipCode) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		weather, err := weatherService.GetWeather(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(weather); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	webServerPort := os.Getenv("WEB_SERVER_PORT")
	log.Printf("Starting web server on port %s", webServerPort)
	http.ListenAndServe(":"+webServerPort, mux)
}
