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

	return env
}
