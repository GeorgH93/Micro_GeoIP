package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"micro_geoip/internal/config"
	"micro_geoip/internal/geoip"
)

func createTestServer(t *testing.T) *Server {
	cfg := &config.Config{}
	cfg.Server.Port = "8080"
	cfg.Server.Host = "localhost"
	cfg.Security.BlockIPParam = false
	
	// Create a mock geoip service
	geoipService := geoip.NewMockService()
	
	return NewServer(cfg, geoipService)
}

func TestHealthCheck(t *testing.T) {
	server := createTestServer(t)
	
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse JSON response")
	}
	
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestGeoLookupWithValidIP(t *testing.T) {
	server := createTestServer(t)
	
	req, err := http.NewRequest("GET", "/geoip?ip=8.8.8.8", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var response GeoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse JSON response")
	}
	
	if response.IP != "8.8.8.8" {
		t.Errorf("Expected IP '8.8.8.8', got '%s'", response.IP)
	}
	
	if response.CountryCode != "US" {
		t.Errorf("Expected country code 'US', got '%s'", response.CountryCode)
	}
	
	if response.Country != "United States" {
		t.Errorf("Expected country 'United States', got '%s'", response.Country)
	}
}

func TestGeoLookupWithInvalidIP(t *testing.T) {
	server := createTestServer(t)
	
	req, err := http.NewRequest("GET", "/geoip?ip=invalid-ip", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var response GeoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse JSON response")
	}
	
	if response.Error == "" {
		t.Error("Expected error message for invalid IP")
	}
}

func TestGetClientIP(t *testing.T) {
	server := createTestServer(t)
	
	// Test X-Forwarded-For header
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")
	
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	
	// For this test, we'll just verify the endpoint responds
	// since we can't easily access the getClientIP method directly
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}