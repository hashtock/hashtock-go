package models

import (
    "net/http"

    "appengine"
    "appengine/datastore"
)

const (
    tagShareKind = "TagShare"
)

type TagShare struct {
    HashTag  string  `json:"hashtag"`
    UserID   string  `json:"user_id"`
    Quantity float64 `json:"quantity"`
}

func (t *TagShare) key(ctx appengine.Context) (key *datastore.Key) {
    return tagShareKey(ctx, t.HashTag, t.UserID)
}

func (t *TagShare) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := t.key(ctx)
    _, err = datastore.Put(ctx, key, t)
    return
}
