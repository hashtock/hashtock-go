package shares

import (
    "appengine"
    "appengine/datastore"
)

func GetOrNew(ctx appengine.Context, hash string, profile_key *datastore.Key, user_id string) (share UserShare) {
    share = UserShare{
        profile_key: profile_key,
        UserID:      user_id,
        HashTag:     hash,
        Quantity:    0,
    }
    key := share.Key(ctx)

    datastore.Get(ctx, key, &share)

    return
}

func Query(ctx appengine.Context) (query *datastore.Query) {
    query = datastore.NewQuery("UserShare").Order("-HashTag")
    return
}
