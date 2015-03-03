package models

import (
    "fmt"
    "net/http"
    "time"

    "appengine"
    "appengine/datastore"
)

func hashTagValueKey(ctx appengine.Context, name string, date time.Time) (key *datastore.Key) {
    return datastore.NewKey(ctx, hashTagValueKind, fmt.Sprintf("%s_%s", name, date.String()), 0, nil)
}

func GetHashTagValues(req *http.Request, hashTagName string) (values []HashTagValue, err error) {
    if err = hashTagExistsOrError(req, hashTagName); err != nil {
        return
    }
    ctx := appengine.NewContext(req)

    since := time.Now().Add(time.Hour * -24)
    q := datastore.NewQuery(hashTagValueKind).Filter("HashTag = ", hashTagName).Filter("Date >", since).Order("Date")
    _, err = q.GetAll(ctx, &values)
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
