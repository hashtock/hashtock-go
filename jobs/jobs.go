package jobs

import (
	// "time"

	"github.com/hashtock/hashtock-go/conf"
	"github.com/hashtock/hashtock-go/models"
)

// Starts all jobs with scheduling from config
// Jobs run forever
func StartJobs(cfg *conf.Config, storage *models.MgoStorage) {
	// bankOrderTicker := time.NewTicker(cfg.Jobs.BankOrders)

	// go func() {
	// 	for range bankOrderTicker.C {
	// 		ExecuteBankOrders()
	// 	}
	// }()
}
