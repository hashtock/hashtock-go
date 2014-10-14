package models

import (
    "appengine"
    "appengine/datastore"
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
