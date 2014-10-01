package profiles

import (
    "appengine"
    "appengine/datastore"

    "shares"
)

type Profile struct {
    UserID string
    Founds float64
}

func Key(ctx appengine.Context, user_id string) (key *datastore.Key) {
    return datastore.NewKey(ctx, "Profile", user_id, 0, nil)
}

func (p *Profile) Key(ctx appengine.Context) (key *datastore.Key) {
    return Key(ctx, p.UserID)
}

func (p *Profile) Put(ctx appengine.Context) (err error) {
    key := p.Key(ctx)
    _, err = datastore.Put(ctx, key, p)
    return
}

func (p *Profile) GetShare(ctx appengine.Context, hash string) (share shares.UserShare) {
    share = shares.GetOrNew(ctx, hash, p.Key(ctx), p.UserID)

    return
}

func (p *Profile) Shares(ctx appengine.Context) (entries []shares.UserShare, err error) {
    userKey := p.Key(ctx)

    q := shares.Query(ctx).Ancestor(userKey)
    _, err = q.GetAll(ctx, &entries)

    return
}
