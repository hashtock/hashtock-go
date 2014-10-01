package shares

import (
    "appengine"
    "appengine/datastore"
)

type UserShare struct {
    UserID   string
    HashTag  string
    Quantity float64

    profile_key *datastore.Key
}

func Key(ctx appengine.Context, hash string, profile_key *datastore.Key) (key *datastore.Key) {
    return datastore.NewKey(ctx, "UserShare", hash, 0, profile_key)
}

func (u *UserShare) Key(ctx appengine.Context) (key *datastore.Key) {
    return Key(ctx, u.UserID, u.profile_key)
}

func (u *UserShare) Put(ctx appengine.Context) (err error) {
    key := u.Key(ctx)
    _, err = datastore.Put(ctx, key, u)
    return
}
