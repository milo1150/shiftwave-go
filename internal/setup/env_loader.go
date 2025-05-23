package setup

import (
	"log"
	"os"
	"shiftwave-go/internal/types"
	"time"

	"github.com/joho/godotenv"
)

func EnvLoader() types.Env {
	env := types.Env{}

	_, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Error loading env")
	}

	// APP ENV
	appEnv := os.Getenv("APP_ENV")
	if err != nil {
		log.Fatalf("Error loading app env: %v", err)
	}
	env.APP_ENV = appEnv

	// Location timezone
	location, err := time.LoadLocation(os.Getenv("LOCAL_TIMEZONE"))
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}
	env.LocalTimezone = location

	// OpenAI
	openAI := os.Getenv("OPENAI_API_KEY")
	if openAI == "" {
		log.Fatalf("openAI token should be not empty")
	}
	env.OpenAI = openAI

	// JWT
	jwt := os.Getenv("JWT")
	if jwt == "" {
		log.Fatalf("jwt secret should be not empty")
	}
	env.JWT = jwt

	// Admin pwd
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatalf("admin password should be not empty")
	}
	env.AdminPassword = adminPassword

	// Redis pwd
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatalf("redis password should be not empty")
	}
	env.RedisPassword = redisPassword

	return env
}
