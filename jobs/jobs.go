package jobs

import (
    "time"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/models"
)

// Starts all jobs with scheduling from config
// Jobs run forever
func StartJobs(cfg *conf.Config, storage *models.MgoStorage) {
    bankOrderTicker := time.NewTicker(cfg.Jobs.BankOrders.Duration())
    tagValueTicker := time.NewTicker(cfg.Jobs.TagValues.Duration())

    go func() {
        for range bankOrderTicker.C {
            ExecuteBankOrders()
        }
    }()

    go func() {
        for range tagValueTicker.C {
            FetchLatestTagValues(cfg)
        }
    }()
}