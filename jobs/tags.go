package jobs

import (
    "log"
    "sync"
    "time"

    "github.com/hashtock/tracker/client"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/models"
)

//TODO(error): Refactor & Tests!
func FetchLatestTagValues(cfg *conf.Config) {
    secret := cfg.Tracker.HMACSecret
    host := cfg.Tracker.Url
    if secret == "" || host == "" {
        log.Println("job:FetchLatestTagValues: Incorrectly configured tracker. Missing values")
        return
    }
    tracker, _ := client.NewTrackerPlain(secret, host)

    latestUpdate, err := models.LatestUpdateToHashTagValues(nil)
    if err != nil {
        log.Println("job:FetchLatestTagValues: Could not figure when last update run. Err:", err)
        return
    }

    tagCountTrend, err := tracker.Trends(latestUpdate.Add(time.Second*1), time.Now())
    if err != nil {
        log.Println("job:FetchLatestTagValues: Failed when fetching data from tracker. Err:", err)
        return
    }

    tags := make(map[string]*models.HashTag, len(tagCountTrend))
    for _, tagCounts := range tagCountTrend {
        tag, err := models.GetOrCreateHashTag(nil, nil, tagCounts.Name)
        if err != nil {
            log.Println("job:FetchLatestTagValues: Could not deal with:", tagCounts.Name, "Because:", err)
            continue
        }
        tags[tag.HashTag] = tag
    }

    for _, tagCounts := range tagCountTrend {
        latestCount := time.Time{}
        latestValue := tags[tagCounts.Name].Value
        for _, count := range tagCounts.Counts {
            tagValue := models.HashTagValue{
                HashTag: tagCounts.Name,
                Value:   float64(count.Count),
                Date:    count.Date,
            }

            if count.Date.After(latestCount) {
                latestValue = float64(count.Count)
            }

            models.AddHashTagValue(nil, tagValue)
        }

        if tags[tagCounts.Name].Value != latestValue {
            tags[tagCounts.Name].Value = latestValue
        } else {
            delete(tags, tagCounts.Name)
        }
    }

    wg := sync.WaitGroup{}
    for _, tag := range tags {
        wg.Add(1)
        go func(innerTag *models.HashTag) {
            defer wg.Done()
            if _, err := models.UpdateHashTagValue(nil, nil, innerTag.HashTag, innerTag.Value); err != nil {
                log.Println("job:FetchLatestTagValues: Could not update:", innerTag.HashTag, "Because:", err)
                return
            }
        }(tag)
    }
    wg.Wait()

    if len(tags) == 0 {
        log.Println("job:FetchLatestTagValues: Tracker has no new tag values yet. Last time updated:", latestUpdate)
    } else {
        log.Printf("job:FetchLatestTagValues: Fetched latest tag values from tracker for %v tags. Last time updated: %v", len(tags), latestUpdate)
    }
}
