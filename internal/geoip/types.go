package geoip

// CountryInfo represents country information from GeoIP lookup
type CountryInfo struct {
	Code string // ISO country code (e.g., "US")
	Name string // Country name (e.g., "United States")
}