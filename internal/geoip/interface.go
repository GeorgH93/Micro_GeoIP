package geoip

// GeoIPService defines the interface for GeoIP lookup services
type GeoIPService interface {
	GetCountry(ip string) (*CountryInfo, error)
	Close() error
}