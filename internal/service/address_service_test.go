package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAddress(t *testing.T) {
	tests := []struct {
		name           string
		zipCode        string
		mockResponse   string
		mockStatusCode int
		expectedResult string
		expectedError  error
	}{
		{
			name:           "valid zipcode",
			zipCode:        "12345678",
			mockResponse:   `{"localidade": "Test City"}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "Test City",
			expectedError:  nil,
		},
		{
			name:           "invalid zipcode format",
			zipCode:        "1234",
			mockResponse:   ``,
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			expectedError:  ErrInvalidZipCode,
		},
		{
			name:           "zipcode not found",
			zipCode:        "87654321",
			mockResponse:   `{"erro": true}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			expectedError:  ErrNotFoundZipCode,
		},
		{
			name:           "server error",
			zipCode:        "12345678",
			mockResponse:   ``,
			mockStatusCode: http.StatusInternalServerError,
			expectedResult: "",
			expectedError:  ErrInvalidZipCode,
		},
		{
			name:           "invalid json response",
			zipCode:        "12345678",
			mockResponse:   `{invalid json}`,
			mockStatusCode: http.StatusOK,
			expectedResult: "",
			expectedError:  errors.New("an error occurred while decoding the response"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			service := NewAddressService(server.URL + "/zipCode")
			result, err := service.GetAddress(tt.zipCode)

			if result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
