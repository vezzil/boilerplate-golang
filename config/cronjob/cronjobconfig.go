package cronjob

import (
	"log"

	"github.com/robfig/cron/v3"
	"boilerplate-golang/config/toml"
)

func InitCronJobs() {
	c := cron.New()

	// Example: cleanup every midnight
	_, err := c.AddFunc(toml.GetConfig().CronJob.CleanupInterval, func() {
		log.Println("🧹 Running cleanup task...")
		// your cleanup logic
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule cleanup job: %v", err)
	}

	// Example: send report every Monday 9 AM
	_, err = c.AddFunc(toml.GetConfig().CronJob.EmailReport, func() {
		log.Println("📧 Sending weekly email report...")
		// your email logic
	})
	if err != nil {
		log.Fatalf("❌ Failed to schedule email job: %v", err)
	}

	c.Start()
	log.Println("✅ Cron jobs initialized")
}
