package models

import (
    "net/http"

    "appengine"
    "appengine/datastore"
    "github.com/gorilla/context"
)

const (
    StartingFounds = 1000
    profileKind    = "Profile"
)

type Profile struct {
    UserID string  `json:"id"`
    Founds float64 `json:"founds"`
}

func (p *Profile) key(ctx appengine.Context) (key *datastore.Key) {
    return profileKey(ctx, p.UserID)
}

func (p *Profile) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := p.key(ctx)
    _, err = datastore.Put(ctx, key, p)

    if err == nil {
        context.Set(req, reqProfile, p)
    }
    return
}
