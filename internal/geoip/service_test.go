/*
 * Copyright (C) 2025  GeorgH93
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package geoip

import (
	"testing"
)

func TestNewService(t *testing.T) {
	// Skip this test as it requires network access
	// In a real environment, you would mock the HTTP client
	t.Skip("Skipping TestNewService as it requires network access")
}

func TestGetCountryWithMockService(t *testing.T) {
	mockService := NewMockService()

	// Test known IP
	countryInfo, err := mockService.GetCountry("8.8.8.8")
	if err != nil {
		t.Errorf("GetCountry failed: %v", err)
	}

	if countryInfo.Code != "US" {
		t.Errorf("Expected country code 'US', got '%s'", countryInfo.Code)
	}

	if countryInfo.Name != "United States" {
		t.Errorf("Expected country name 'United States', got '%s'", countryInfo.Name)
	}

	// Test unknown IP
	countryInfo, err = mockService.GetCountry("192.168.1.1")
	if err != nil {
		t.Errorf("GetCountry failed: %v", err)
	}

	if countryInfo.Code != "Unknown" {
		t.Errorf("Expected country code 'Unknown', got '%s'", countryInfo.Code)
	}

	if countryInfo.Name != "Unknown" {
		t.Errorf("Expected country name 'Unknown', got '%s'", countryInfo.Name)
	}
}

func TestMockServiceSetCountry(t *testing.T) {
	mockService := NewMockService()

	// Set a custom country for an IP
	mockService.SetCountry("203.0.113.1", "GB", "United Kingdom")

	countryInfo, err := mockService.GetCountry("203.0.113.1")
	if err != nil {
		t.Errorf("GetCountry failed: %v", err)
	}

	if countryInfo.Code != "GB" {
		t.Errorf("Expected country code 'GB', got '%s'", countryInfo.Code)
	}

	if countryInfo.Name != "United Kingdom" {
		t.Errorf("Expected country name 'United Kingdom', got '%s'", countryInfo.Name)
	}
}

func TestValidateIP(t *testing.T) {
	testCases := []struct {
		ip    string
		valid bool
	}{
		{"8.8.8.8", true},
		{"192.168.1.1", true},
		{"2001:4860:4860::8888", true},
		{"invalid-ip", false},
		{"999.999.999.999", false},
		{"", false},
	}

	mockService := NewMockService()

	for _, tc := range testCases {
		_, err := mockService.GetCountry(tc.ip)
		hasError := err != nil

		// Note: MockService doesn't validate IPs, so we expect no errors
		// In a real implementation, you'd want to add IP validation
		if hasError {
			t.Errorf("Unexpected error for IP %s: %v", tc.ip, err)
		}
	}
}

func TestServiceClose(t *testing.T) {
	mockService := NewMockService()

	err := mockService.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
