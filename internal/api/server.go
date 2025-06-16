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

package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"micro_geoip/internal/config"
	"micro_geoip/internal/geoip"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config       *config.Config
	geoipService geoip.GeoIPService
	router       *gin.Engine
}

type GeoResponse struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Error       string `json:"error,omitempty"`
}

func NewServer(cfg *config.Config, geoipService geoip.GeoIPService) *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		config:       cfg,
		geoipService: geoipService,
		router:       gin.New(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Add basic middleware
	s.router.Use(gin.Recovery())
	s.router.Use(gin.Logger())

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// GeoIP lookup endpoints
	s.router.GET("/", s.geoLookup)
	s.router.GET("/geoip", s.geoLookup)
	s.router.GET("/geoip/:ip", s.geoLookupWithIP)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("Starting server on %s\n", addr)
	return s.router.Run(addr)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) geoLookup(c *gin.Context) {
	var targetIP string

	// Check if IP parameter is blocked
	if s.config.Security.BlockIPParam {
		targetIP = s.getClientIP(c)
	} else {
		// Try to get IP from query parameter
		if ip := c.Query("ip"); ip != "" {
			targetIP = ip
		} else {
			targetIP = s.getClientIP(c)
		}
	}

	s.performGeoLookup(c, targetIP)
}

func (s *Server) geoLookupWithIP(c *gin.Context) {
	// If IP parameter is blocked, ignore the path parameter and use client IP
	if s.config.Security.BlockIPParam {
		targetIP := s.getClientIP(c)
		s.performGeoLookup(c, targetIP)
		return
	}

	ip := c.Param("ip")
	s.performGeoLookup(c, ip)
}

func (s *Server) performGeoLookup(c *gin.Context, ip string) {
	// Validate IP address
	if net.ParseIP(ip) == nil {
		c.JSON(http.StatusBadRequest, GeoResponse{
			IP:    ip,
			Error: "Invalid IP address",
		})
		return
	}

	// Perform GeoIP lookup
	countryInfo, err := s.geoipService.GetCountry(ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeoResponse{
			IP:    ip,
			Error: fmt.Sprintf("GeoIP lookup failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, GeoResponse{
		IP:          ip,
		Country:     countryInfo.Name,
		CountryCode: countryInfo.Code,
	})
}

func (s *Server) getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}
