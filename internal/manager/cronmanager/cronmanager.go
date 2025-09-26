package cronmanager

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"boilerplate-golang/internal/config"
)

var c *cron.Cron

// Init starts the cron scheduler and registers jobs from config.
// If no cron jobs are configured, the scheduler won't be started.
func Init() {
	cfg := config.Get()

	// Skip initialization if no cron jobs are configured
	if cfg.CronJob.CleanupInterval == "" && cfg.CronJob.EmailReport == "" {
		log.Println("cron: no cron jobs configured, skipping cron scheduler")
		return
	}

	loc := time.Local // could be made configurable later
	c = cron.New(cron.WithLocation(loc))

	jobsScheduled := 0

	if cfg.CronJob.CleanupInterval != "" {
		if _, err := c.AddFunc(cfg.CronJob.CleanupInterval, func() {
			log.Println("cron: running cleanup task")
			// TODO: add cleanup logic or call into a job package
		}); err != nil {
			log.Printf("cron: failed to schedule cleanup: %v", err)
		} else {
			jobsScheduled++
		}
	}

	if cfg.CronJob.EmailReport != "" {
		if _, err := c.AddFunc(cfg.CronJob.EmailReport, func() {
			log.Println("cron: sending email report")
			// TODO: add email report logic
		}); err != nil {
			log.Printf("cron: failed to schedule email report: %v", err)
		} else {
			jobsScheduled++
		}
	}

	if jobsScheduled > 0 {
		log.Printf("cron: started with %d job(s) scheduled", jobsScheduled)
		c.Start()
	} else {
		log.Println("cron: no valid cron jobs scheduled, not starting cron scheduler")
	}
	log.Println("cronmanager: started")
}

// Stop halts the cron scheduler.
func Stop() {
	if c != nil {
		c.Stop()
	}
}
