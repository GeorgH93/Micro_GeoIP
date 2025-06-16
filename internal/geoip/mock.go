package geoip

// MockService implements the GeoIP service interface for testing
type MockService struct {
	CountryMap map[string]*CountryInfo
}

func NewMockService() *MockService {
	return &MockService{
		CountryMap: map[string]*CountryInfo{
			"8.8.8.8":                 {Code: "US", Name: "United States"},
			"1.1.1.1":                 {Code: "US", Name: "United States"},
			"208.67.222.222":          {Code: "US", Name: "United States"},
			"134.195.196.26":          {Code: "DE", Name: "Germany"},
			"2001:4860:4860::8888":    {Code: "US", Name: "United States"},
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