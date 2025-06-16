package config

import (
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" env:"PORT"`
		Host string `yaml:"host" env:"HOST"`
	} `yaml:"server"`
	
	GeoIP struct {
		MaxMindAPIKey   string `yaml:"maxmind_api_key" env:"MAXMIND_API_KEY"`
		DatabasePath    string `yaml:"database_path" env:"GEOIP_DB_PATH"`
		UpdateInterval  string `yaml:"update_interval" env:"GEOIP_UPDATE_INTERVAL"`
		MaxMindURL      string `yaml:"maxmind_url" env:"MAXMIND_DOWNLOAD_URL"`
		DBIPUrl         string `yaml:"dbip_url" env:"DBIP_DOWNLOAD_URL"`
		PreferDBIP      bool   `yaml:"prefer_dbip" env:"PREFER_DBIP"`
	} `yaml:"geoip"`
	
	Security struct {
		BlockIPParam bool `yaml:"block_ip_param" env:"BLOCK_IP_PARAM"`
	} `yaml:"security"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	
	// Set defaults
	cfg.Server.Port = "8080"
	cfg.Server.Host = "0.0.0.0"
	cfg.GeoIP.DatabasePath = "./data/GeoLite2-Country.mmdb"
	cfg.GeoIP.UpdateInterval = "720h" // 30 days
	cfg.GeoIP.MaxMindURL = "https://download.maxmind.com/app/geoip_download"
	cfg.GeoIP.DBIPUrl = "https://download.db-ip.com/free/dbip-country-lite-{YYYY-MM}.mmdb.gz"
	cfg.GeoIP.PreferDBIP = false
	cfg.Security.BlockIPParam = false
	
	// Try to load from config file
	if err := loadFromFile(cfg); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	
	// Override with environment variables
	loadFromEnv(cfg)
	
	return cfg, nil
}

func loadFromFile(cfg *Config) error {
	configPaths := []string{
		"config.yaml",
		"config.yml",
		"/etc/micro_geoip/config.yaml",
		"/etc/micro_geoip/config.yml",
	}
	
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return yaml.Unmarshal(data, cfg)
		}
	}
	
	return os.ErrNotExist
}

func loadFromEnv(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		cfg.Server.Host = host
	}
	if apiKey := os.Getenv("MAXMIND_API_KEY"); apiKey != "" {
		cfg.GeoIP.MaxMindAPIKey = apiKey
	}
	if dbPath := os.Getenv("GEOIP_DB_PATH"); dbPath != "" {
		cfg.GeoIP.DatabasePath = dbPath
	}
	if updateInterval := os.Getenv("GEOIP_UPDATE_INTERVAL"); updateInterval != "" {
		cfg.GeoIP.UpdateInterval = updateInterval
	}
	if maxmindURL := os.Getenv("MAXMIND_DOWNLOAD_URL"); maxmindURL != "" {
		cfg.GeoIP.MaxMindURL = maxmindURL
	}
	if dbipURL := os.Getenv("DBIP_DOWNLOAD_URL"); dbipURL != "" {
		cfg.GeoIP.DBIPUrl = dbipURL
	}
	if preferDBIP := os.Getenv("PREFER_DBIP"); preferDBIP != "" {
		if val, err := strconv.ParseBool(preferDBIP); err == nil {
			cfg.GeoIP.PreferDBIP = val
		}
	}
	if blockIP := os.Getenv("BLOCK_IP_PARAM"); blockIP != "" {
		if val, err := strconv.ParseBool(blockIP); err == nil {
			cfg.Security.BlockIPParam = val
		}
	}
}

func (c *Config) GetDatabaseDir() string {
	return filepath.Dir(c.GeoIP.DatabasePath)
}