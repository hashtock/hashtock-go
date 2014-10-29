package models

import (
    "net/http"

    "appengine"
    "appengine/datastore"
)

const (
    hashTagKind = "HashTag"
)

type HashTag struct {
    HashTag string  `json:"hashtag"`
    Value   float64 `json:"value"`
    InBank  float64 `json:"in_bank"`
}

func (h *HashTag) key(ctx appengine.Context) (key *datastore.Key) {
    return hashTagKey(ctx, h.HashTag)
}

func (h *HashTag) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := h.key(ctx)
    _, err = datastore.Put(ctx, key, h)
    return
}
