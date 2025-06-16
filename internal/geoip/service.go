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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"micro_geoip/internal/config"

	"github.com/oschwald/geoip2-golang"
	"github.com/robfig/cron/v3"
)

type Service struct {
	config *config.Config
	db     *geoip2.Reader
	cron   *cron.Cron
}

func NewService(cfg *config.Config) (*Service, error) {
	s := &Service{
		config: cfg,
		cron:   cron.New(),
	}

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.GetDatabaseDir(), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Try to load existing database
	if err := s.loadDatabase(); err != nil {
		log.Printf("Failed to load existing database: %v", err)

		// If no database exists, download it
		log.Println("Downloading initial GeoIP database...")
		if err := s.downloadDatabase(); err != nil {
			return nil, fmt.Errorf("failed to download initial database: %w", err)
		}

		// Try to load the downloaded database
		if err := s.loadDatabase(); err != nil {
			return nil, fmt.Errorf("failed to load downloaded database: %w", err)
		}
	}

	// Set up automatic updates
	s.setupAutoUpdate()

	return s, nil
}

func (s *Service) loadDatabase() error {
	if _, err := os.Stat(s.config.GeoIP.DatabasePath); os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist: %s", s.config.GeoIP.DatabasePath)
	}

	db, err := geoip2.Open(s.config.GeoIP.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open GeoIP database: %w", err)
	}

	// Close old database if exists
	if s.db != nil {
		s.db.Close()
	}

	s.db = db
	log.Printf("GeoIP database loaded: %s", s.config.GeoIP.DatabasePath)
	return nil
}

func (s *Service) downloadDatabase() error {
	// Try MaxMind first if API key is available and not preferring DB-IP
	if s.config.GeoIP.MaxMindAPIKey != "" && !s.config.GeoIP.PreferDBIP {
		if err := s.downloadMaxMindDatabase(); err != nil {
			log.Printf("MaxMind download failed: %v, trying DB-IP fallback", err)
			return s.downloadDBIPDatabase()
		}
		return nil
	}

	// Try DB-IP first (free database)
	if err := s.downloadDBIPDatabase(); err != nil {
		log.Printf("DB-IP download failed: %v", err)

		// Fallback to MaxMind if API key is available
		if s.config.GeoIP.MaxMindAPIKey != "" {
			log.Printf("Trying MaxMind fallback...")
			return s.downloadMaxMindDatabase()
		}

		return fmt.Errorf("no database source available: %w", err)
	}

	return nil
}

func (s *Service) downloadMaxMindDatabase() error {
	if s.config.GeoIP.MaxMindAPIKey == "" {
		return fmt.Errorf("no MaxMind API key provided")
	}

	// Build download URL
	downloadURL := fmt.Sprintf("%s?edition_id=GeoLite2-Country&license_key=%s&suffix=tar.gz",
		s.config.GeoIP.MaxMindURL,
		url.QueryEscape(s.config.GeoIP.MaxMindAPIKey))

	log.Printf("Downloading GeoIP database from MaxMind...")

	// Download the tar.gz file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download from MaxMind: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MaxMind download failed with status: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "maxmind-geoip-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Save downloaded content
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Extract the database file
	if err := s.extractMaxMindDatabase(tmpFile.Name()); err != nil {
		return fmt.Errorf("failed to extract MaxMind database: %w", err)
	}

	log.Printf("MaxMind GeoIP database downloaded and extracted successfully")
	return nil
}

func (s *Service) downloadDBIPDatabase() error {
	// Format current date for DB-IP URL (YYYY-MM format)
	currentDate := time.Now().Format("2006-01")
	downloadURL := strings.Replace(s.config.GeoIP.DBIPUrl, "{YYYY-MM}", currentDate, 1)

	log.Printf("Downloading GeoIP database from DB-IP...")

	// Download the .mmdb.gz file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download from DB-IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DB-IP download failed with status: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "dbip-geoip-*.mmdb.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Save downloaded content
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Extract the database file
	if err := s.extractDBIPDatabase(tmpFile.Name()); err != nil {
		return fmt.Errorf("failed to extract DB-IP database: %w", err)
	}

	log.Printf("DB-IP GeoIP database downloaded and extracted successfully")
	return nil
}

func (s *Service) extractMaxMindDatabase(tarGzPath string) error {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Look for the .mmdb file
		if strings.HasSuffix(header.Name, ".mmdb") && strings.Contains(header.Name, "GeoLite2-Country") {
			// Extract to the configured path
			outFile, err := os.Create(s.config.GeoIP.DatabasePath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("GeoLite2-Country.mmdb not found in MaxMind archive")
}

func (s *Service) extractDBIPDatabase(gzPath string) error {
	file, err := os.Open(gzPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// Create output file
	outFile, err := os.Create(s.config.GeoIP.DatabasePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy decompressed content
	if _, err := io.Copy(outFile, gzr); err != nil {
		return err
	}

	return nil
}

func (s *Service) setupAutoUpdate() {
	// Parse update interval
	_, err := time.ParseDuration(s.config.GeoIP.UpdateInterval)
	if err != nil {
		log.Printf("Invalid update interval '%s', using default (30 days): %v", s.config.GeoIP.UpdateInterval, err)
	}

	// Convert to cron format (approximate - runs once per month)
	cronSpec := "0 0 1 * *" // First day of every month at midnight

	_, err = s.cron.AddFunc(cronSpec, func() {
		log.Println("Starting scheduled GeoIP database update...")
		if err := s.downloadDatabase(); err != nil {
			log.Printf("Scheduled database update failed: %v", err)
			return
		}

		if err := s.loadDatabase(); err != nil {
			log.Printf("Failed to reload database after update: %v", err)
			return
		}

		log.Println("Scheduled GeoIP database update completed successfully")
	})

	if err != nil {
		log.Printf("Failed to schedule database updates: %v", err)
		return
	}

	s.cron.Start()
	log.Printf("Scheduled automatic database updates: %s", cronSpec)
}

func (s *Service) GetCountry(ip string) (*CountryInfo, error) {
	if s.db == nil {
		return nil, fmt.Errorf("GeoIP database not available")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	record, err := s.db.Country(parsedIP)
	if err != nil {
		return nil, fmt.Errorf("GeoIP lookup failed: %w", err)
	}

	countryInfo := &CountryInfo{
		Code: "Unknown",
		Name: "Unknown",
	}

	// Set country code
	if record.Country.IsoCode != "" {
		countryInfo.Code = record.Country.IsoCode
	}

	// Set country name (prefer English)
	if record.Country.Names["en"] != "" {
		countryInfo.Name = record.Country.Names["en"]
	} else if len(record.Country.Names) > 0 {
		// Fallback to any available name
		for _, name := range record.Country.Names {
			countryInfo.Name = name
			break
		}
	}

	return countryInfo, nil
}

func (s *Service) Close() error {
	if s.cron != nil {
		s.cron.Stop()
	}

	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
