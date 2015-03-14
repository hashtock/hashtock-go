package models

import (
    "fmt"
    "net/http"
    "sort"
    "time"

    "appengine"
    "appengine/datastore"
)

func hashTagValueKey(ctx appengine.Context, name string, date time.Time) (key *datastore.Key) {
    return datastore.NewKey(ctx, hashTagValueKind, fmt.Sprintf("%s_%s", name, date.String()), 0, nil)
}

func GetHashTagValues(req *http.Request, hashTagName string, since time.Time, sampling time.Duration) (values []HashTagValue, err error) {
    if err = hashTagExistsOrError(req, hashTagName); err != nil {
        return
    }
    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(hashTagValueKind).Filter("HashTag = ", hashTagName).Filter("Date >", since).Order("Date")
    _, err = q.GetAll(ctx, &values)
    if err != nil {
        return
    }

    if sampling != 0 {
        return resampleTagValues(values, sampling), nil
    }

    return
}

func LatestUpdateToHashTagValues(req *http.Request) (date time.Time, err error) {
    ctx := appengine.NewContext(req)
    q := datastore.NewQuery(hashTagValueKind).Order("-Date").Limit(1)

    values := []HashTagValue{}
    _, err = q.GetAll(ctx, &values)
    if len(values) == 0 {
        return time.Time{}, err
    }

    return values[0].Date, err
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
