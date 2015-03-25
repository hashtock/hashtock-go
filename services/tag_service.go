package services

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/hashtock-go/core"
    "github.com/hashtock/hashtock-go/models"
)

const (
    lastDay      = "1"
    lastWeek     = "7"
    lastTwoWeeks = "14"
    lastMonth    = "30"
)

var (
    showLast = map[string]struct {
        sample   time.Duration
        duration time.Duration
    }{
        lastDay:      {0, 24 * time.Hour},
        lastWeek:     {24 * time.Hour, 7 * 24 * time.Hour},
        lastTwoWeeks: {24 * time.Hour, 14 * 24 * time.Hour},
        lastMonth:    {24 * time.Hour, 30 * 24 * time.Hour},
    }
)

// List of all tags with bank values
func ListOfAllHashTags(req *http.Request, r render.Render) {
    tags, err := models.GetAllHashTags(req)

    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tags)
}

// Details about the hash tag
func TagInfo(req *http.Request, params martini.Params, r render.Render) {
    hash_tag_name := params["tag"]

    tag, err := models.GetHashTag(req, hash_tag_name)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tag)
}

// Add new tag (admin)
func NewHashTag(req *http.Request, r render.Render) {
    tag := models.HashTag{}

    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&tag); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    new_tag, err := models.AddHashTag(req, profile, tag)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusCreated, new_tag)
}

// Set tag value (admin)
func SetTagValue(req *http.Request, params martini.Params, r render.Render) {
    hash_tag_name := params["tag"]

    updated_tag := models.HashTag{}
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&updated_tag); err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    if updated_tag.HashTag != "" && hash_tag_name != updated_tag.HashTag {
        err := core.NewBadRequestError("hashtag value has to be empty or correct")
        r.JSON(core.ErrToErrorer(err))
        return
    }

    profile, err := models.GetProfile(req)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    tag, err := models.UpdateHashTagValue(req, profile, hash_tag_name, updated_tag.Value)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    r.JSON(http.StatusOK, tag)
}

func TagValues(req *http.Request, params martini.Params, r render.Render) {
    hashTagName := params["tag"]

    queryValues := req.URL.Query()
    days := lastDay
    if daysStr := queryValues.Get("days"); daysStr != "" {
        days = daysStr
    }

    def, ok := showLast[days]
    if !ok {
        herr := core.NewBadRequestError(fmt.Sprintf("Duration %s not supported", days))
        r.JSON(core.ErrToErrorer(herr))
        return
    }

    since := time.Now().Add(-def.duration)
    tagValues, err := models.GetHashTagValues(req, hashTagName, since, def.sample)
    if err != nil {
        r.JSON(core.ErrToErrorer(err))
        return
    }

    // Subtract 1h to make room for sampling error
    if len(tagValues) > 0 && def.sample != 0 && time.Since(tagValues[0].Date) < def.duration {
        tagValues = tagValues[:0]
    }

    r.JSON(http.StatusOK, tagValues)
}
