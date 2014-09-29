package profiles

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
)

func CurrentProfile(ctx appengine.Context) (profile Profile, err error) {
    u := user.Current(ctx)
    key := Key(ctx, u.Email)

    err = datastore.Get(ctx, key, &profile)
    return
}

func CreateNewUserIfDoesNotExist(ctx appengine.Context, current_user user.User) (err error) {
    key := Key(ctx, current_user.Email)

    profile := Profile{
        UserID: current_user.Email,
        Founds: 1000,
    }

    if datastore.Get(ctx, key, &profile) != datastore.ErrNoSuchEntity {
        return
    }

    _, err = datastore.Put(ctx, key, &profile)
    return
}
