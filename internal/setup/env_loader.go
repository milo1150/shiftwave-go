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
	env.OpenAI = openAI

	// JWT
	jwt := os.Getenv("JWT")
	env.JWT = jwt

	return env
}
