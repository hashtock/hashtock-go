package models

import (
    "net/http"
    "time"

    "appengine"
    "appengine/datastore"
)

const (
    hashTagValueKind = "HashTagValue"
)

type HashTagValue struct {
    HashTag string    `json:"-"`
    Value   float64   `json:"value"`
    Date    time.Time `json:"date"`
}

func (h *HashTagValue) key(ctx appengine.Context) (key *datastore.Key) {
    return hashTagValueKey(ctx, h.HashTag, h.Date)
}

func (h *HashTagValue) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := h.key(ctx)
    _, err = datastore.Put(ctx, key, h)
    return
}
