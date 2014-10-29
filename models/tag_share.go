package models

import (
    "fmt"
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
    uniq_id := fmt.Sprintf("%s-%s", t.HashTag, t.UserID)
    return tagShareKey(ctx, uniq_id)
}

func (t *TagShare) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := t.key(ctx)
    _, err = datastore.Put(ctx, key, t)
    return
}
