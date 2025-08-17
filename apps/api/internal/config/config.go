package config

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	IsProd       bool
	CookieDomain string
	CookieMaxAge time.Duration
	CookieSecure bool
}

func Load() Config {
	isProd := gin.Mode() == gin.ReleaseMode

	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		if isProd {
			log.Fatalf("COOKIE_DOMAIN env var is not set")

		} else {
			domain = "localhost"
		}
	}

	return Config{
		IsProd:       isProd,
		CookieDomain: domain,
		CookieSecure: isProd,
		CookieMaxAge: time.Hour * 24,
	}
}
