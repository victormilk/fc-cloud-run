package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherService_GetWeather(t *testing.T) {
	tests := []struct {
		name       string
		city       string
		apiKey     string
		apiURL     string
		mockResp   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful response",
			city:       "London",
			apiKey:     "test_api_key",
			apiURL:     "http://example.com",
			mockResp:   `{"current": {"temp_c": 20.0, "temp_f": 68.0}}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid API key",
			city:       "London",
			apiKey:     "invalid_api_key",
			apiURL:     "http://example.com",
			mockResp:   `{"error": "invalid API key"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "invalid city",
			city:       "InvalidCity",
			apiKey:     "test_api_key",
			apiURL:     "http://example.com",
			mockResp:   `{"error": "invalid city"}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "malformed JSON response",
			city:       "London",
			apiKey:     "test_api_key",
			apiURL:     "http://example.com",
			mockResp:   `{"current": {"temp_c": "invalid", "temp_f": 68.0}}`,
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.mockResp))
			}))
			defer server.Close()

			ws := NewWeatherService(tt.apiKey, server.URL)
			_, err := ws.GetWeather(tt.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWeather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
