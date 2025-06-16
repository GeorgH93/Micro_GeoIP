package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("GEOIP_DB_PATH", "/tmp/test-geoip.mmdb")
	os.Setenv("PORT", "8081")
	os.Setenv("MAXMIND_API_KEY", "test-key")
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	os.Remove("/tmp/test-geoip.mmdb")
	os.Unsetenv("GEOIP_DB_PATH")
	os.Unsetenv("PORT")
	os.Unsetenv("MAXMIND_API_KEY")
	
	os.Exit(code)
}

func TestImports(t *testing.T) {
	// This test ensures that all imports are valid
	// The fact that this file compiles means the imports work
	t.Log("All imports are valid")
}