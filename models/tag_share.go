package models

import (
    "net/http"

    "appengine"
    "appengine/datastore"
)

const (
    tagShareKind = "TagShare"
    minShareStep = 0.01
)

type TagShare struct {
    HashTag  string  `json:"hashtag"`
    UserID   string  `json:"-"`
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

func (t *TagShare) Delete(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := t.key(ctx)
    err = datastore.Delete(ctx, key)
    return
}
