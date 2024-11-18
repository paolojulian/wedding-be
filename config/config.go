package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	appURI       string
	adminURI     string
	cookieDomain string
)

func init() {
	// Load .env file only once when the package is imported
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Cache environment variables
	appURI = os.Getenv("APP_URI")
	if appURI == "" {
		log.Fatal("APP_URI environment variable not set")
	}

	adminURI = os.Getenv("ADMIN_URI")
	if adminURI == "" {
		log.Fatal("ADMIN_URI environment variable not set")
	}

	cookieDomain = os.Getenv("ADMIN_URI")
	// Note that cookie domain can be ""
	// if cookieDomain == "" {
	// 	log.Fatal("COOKIE_DOMAIN environment variable not set")
	// }
}

// GetAppURI returns the cached APP_URI environment variable
func GetAppURI() string {
	return appURI
}

// GetAdminURI returns the cached ADMIN_URI environment variable
func GetAdminURI() string {
	return adminURI
}

func GetCookieDomain() string {
	return cookieDomain
}
