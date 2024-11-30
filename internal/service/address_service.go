package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")
	ErrNotFoundZipCode = errors.New("can not find zipcode")
)

type AddressService struct {
	URL string
}

func NewAddressService(url string) *AddressService {
	return &AddressService{
		URL: url,
	}
}

func (s *AddressService) GetAddress(zipCode string) (string, error) {
	if err := validateZipCode(zipCode); err != nil {
		return "", err
	}
	url := strings.Replace(s.URL, "zipCode", zipCode, 1)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Println("failed to get address", err)
		return "", errors.New("failed to get address")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrInvalidZipCode
	}

	var apiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", errors.New("an error occurred while decoding the response")
	}

	if _, ok := apiResp["erro"]; ok {
		return "", ErrNotFoundZipCode
	}

	return apiResp["localidade"].(string), nil
}

func validateZipCode(zipCode string) error {
	regex := regexp.MustCompile(`^\d{8}$`)
	if !regex.MatchString(zipCode) {
		return ErrInvalidZipCode
	}
	return nil
}
