package services

import (
    "log"
    "net/http"
    "sync"
    "time"

    "appengine"
    "appengine/urlfetch"
    "github.com/hashtock/tracker/client"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/conf"
    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

// Trigger execution of bank orders
// TODO(access): This endpoind should only by available to admin or cron
func ExecuteBankOrders(req *http.Request, r render.Render) {
    activeOrders, err := models.GetAllActiveBankOrders(req)
    if err != nil {
        r.Error(http.StatusInternalServerError)
        return
    }

    //TODO(error): Handle errors somehow
    for _, order := range activeOrders {
        if err := models.ExecuteBankOrder(req, order); err != nil {
            log.Println(err)
        }
    }

    r.Status(http.StatusOK)
}

//TODO(error): Refactor & Tests!
func FetchLatestTagValues(req *http.Request, r render.Render) {
    c := appengine.NewContext(req)
    secret, host, err := conf.TrackerSecretAndHost(req)

    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    } else if secret == "" || host == "" {
        err = core.NewInternalServerError("Tracker auth issues")
        r.JSON(core.ErrToErrorer(err))
        return
    }

    tracker, _ := client.NewTrackerPlain(secret, host)
    tracker.Client = urlfetch.Client(c)

    latestUpdate, err := models.LatestUpdateToHashTagValues(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    tagCountTrend, err := tracker.Trends(latestUpdate.Add(time.Second*1), time.Now())
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    tags := make(map[string]*models.HashTag, len(tagCountTrend))
    for _, tagCounts := range tagCountTrend {
        tag, err := models.GetOrCreateHashTag(req, tagCounts.Name)
        if err != nil {
            log.Println("Could not deal with:", tagCounts.Name, "Because:", err)
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

            tagValue.Put(req)
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
            // ToDo: It would be safer to update only 1 field here
            if _, err := models.UpdateHashTagValue(req, innerTag.HashTag, innerTag.Value); err != nil {
                log.Println("Could not update:", innerTag.HashTag, "Because:", err)
                return
            }
        }(tag)
    }
    wg.Wait()

    r.Status(http.StatusOK)
}
