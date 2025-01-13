package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func OpenAITranslateScheduler() {
	// Create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Error init scheduler: %v", err)
	}

	count := 0

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				count++
				log.Println(count)

				if count == 10 {
					s.Shutdown()
					s.StopJobs()
				}
			},
		),
	)
	if err != nil {
		log.Fatalf("Error init NewJob: %v", err)
	}

	// Start Cron
	s.Start()
}
