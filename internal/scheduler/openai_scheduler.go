package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"

	"shiftwave-go/internal/services"
	"shiftwave-go/internal/types"
)

func InitializeOpenAiTranslateScheduler(app *types.App) {
	// Create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error init scheduler: %v", err)
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				services.TranslateAndUpdateMyanmarReviews(app)
			},
		),
	)
	if err != nil {
		log.Fatalf("Error init NewJob: %v", err)
	}

	// Start Cron
	s.Start()
}
