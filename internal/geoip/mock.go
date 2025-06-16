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

// MockService implements the GeoIP service interface for testing
type MockService struct {
	CountryMap map[string]*CountryInfo
}

func NewMockService() *MockService {
	return &MockService{
		CountryMap: map[string]*CountryInfo{
			"8.8.8.8":              {Code: "US", Name: "United States"},
			"1.1.1.1":              {Code: "US", Name: "United States"},
			"208.67.222.222":       {Code: "US", Name: "United States"},
			"134.195.196.26":       {Code: "DE", Name: "Germany"},
			"2001:4860:4860::8888": {Code: "US", Name: "United States"},
		},
	}
}

func (m *MockService) GetCountry(ip string) (*CountryInfo, error) {
	if m.CountryMap == nil {
		m.CountryMap = make(map[string]*CountryInfo)
		m.CountryMap["8.8.8.8"] = &CountryInfo{Code: "US", Name: "United States"}
		m.CountryMap["1.1.1.1"] = &CountryInfo{Code: "US", Name: "United States"}
	}

	if country, exists := m.CountryMap[ip]; exists {
		return country, nil
	}

	// Default response for unknown IPs
	return &CountryInfo{Code: "Unknown", Name: "Unknown"}, nil
}

func (m *MockService) Close() error {
	return nil
}

func (m *MockService) SetCountry(ip, code, name string) {
	if m.CountryMap == nil {
		m.CountryMap = make(map[string]*CountryInfo)
	}
	m.CountryMap[ip] = &CountryInfo{Code: code, Name: name}
}

func (m *MockService) AddError(ip string) {
	if m.CountryMap == nil {
		m.CountryMap = make(map[string]*CountryInfo)
	}
	m.CountryMap[ip] = &CountryInfo{Code: "ERROR", Name: "ERROR"}
}
