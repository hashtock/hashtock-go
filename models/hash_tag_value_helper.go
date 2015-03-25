package models

import (
    "net/http"
    "sort"
    "time"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func GetHashTagValues(req *http.Request, hashTagName string, since time.Time, sampling time.Duration) (values []HashTagValue, err error) {
    if err = hashTagExistsOrError(req, hashTagName); err != nil {
        return
    }

    col := storage.Collection(HashTagValueCollectionName)
    defer col.Database.Session.Close()

    since = since.Truncate(24 * time.Hour)
    selector := bson.M{
        "hashtag": hashTagName,
        "date": bson.M{
            "$gt": since,
        },
    }

    err = col.Find(selector).Sort("date").All(&values)
    if err != nil {
        return
    }

    if sampling != 0 {
        return resampleTagValues(values, sampling), nil
    }

    return
}

func LatestUpdateToHashTagValues(req *http.Request) (date time.Time, err error) {
    col := storage.Collection(HashTagValueCollectionName)
    defer col.Database.Session.Close()

    value := HashTagValue{}
    err = col.Find(nil).Sort("-date").Limit(1).One(&value)

    if err == mgo.ErrNotFound {
        return time.Time{}, nil
    }

    return value.Date, err
}

func AddHashTagValue(req *http.Request, value HashTagValue) (err error) {
    col := storage.Collection(HashTagValueCollectionName)
    defer col.Database.Session.Close()

    err = col.Insert(value)
    return
}

func resampleTagValues(values []HashTagValue, sampling time.Duration) (sorted []HashTagValue) {
    if sampling == 0 || len(values) == 0 {
        return values
    }
    tag := values[0].HashTag

    type sampleMap map[time.Time]struct {
        count float64
        value float64
    }
    mapped := make(sampleMap, 0)

    for _, value := range values {
        resampledDate := value.Date.Truncate(sampling)
        currentValue := mapped[resampledDate]
        currentValue.value += value.Value
        currentValue.count += 1
        mapped[resampledDate] = currentValue
    }

    sorted = make([]HashTagValue, 0, len(mapped))
    for resampledDate, aggregatedValues := range mapped {
        resampledValue := aggregatedValues.value / aggregatedValues.count
        sorted = append(sorted, HashTagValue{tag, resampledValue, resampledDate})
    }

    sort.Sort(byDate(sorted))

    return
}
