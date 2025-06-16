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

package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}

	if cfg.Security.BlockIPParam != false {
		t.Errorf("Expected default BlockIPParam false, got %v", cfg.Security.BlockIPParam)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("MAXMIND_API_KEY", "test-key")
	os.Setenv("BLOCK_IP_PARAM", "true")
	os.Setenv("PREFER_DBIP", "true")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("MAXMIND_API_KEY")
		os.Unsetenv("BLOCK_IP_PARAM")
		os.Unsetenv("PREFER_DBIP")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port from env 9090, got %s", cfg.Server.Port)
	}

	if cfg.GeoIP.MaxMindAPIKey != "test-key" {
		t.Errorf("Expected MaxMind API key from env test-key, got %s", cfg.GeoIP.MaxMindAPIKey)
	}

	if cfg.Security.BlockIPParam != true {
		t.Errorf("Expected BlockIPParam from env true, got %v", cfg.Security.BlockIPParam)
	}

	if cfg.GeoIP.PreferDBIP != true {
		t.Errorf("Expected PreferDBIP from env true, got %v", cfg.GeoIP.PreferDBIP)
	}
}

func TestGetDatabaseDir(t *testing.T) {
	cfg := &Config{}
	cfg.GeoIP.DatabasePath = "/tmp/data/GeoLite2-Country.mmdb"

	expected := "/tmp/data"
	if dir := cfg.GetDatabaseDir(); dir != expected {
		t.Errorf("Expected database dir %s, got %s", expected, dir)
	}
}
