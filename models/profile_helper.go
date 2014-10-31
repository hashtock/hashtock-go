package models

import (
    "errors"
    "net/http"

    "appengine"
    "appengine/datastore"
    "appengine/user"
    "github.com/gorilla/context"
)

const (
    reqProfile = "reqProfile"
)

func userId(u *user.User) string {
    return u.Email
}

func profileKey(ctx appengine.Context, user_id string) (key *datastore.Key) {
    return datastore.NewKey(ctx, profileKind, user_id, 0, nil)
}

func profileKeyForUser(ctx appengine.Context, u *user.User) *datastore.Key {
    return profileKey(ctx, userId(u))
}

// Creates new Profile for given User. It does not save it to DB
func newUserProfile(u *user.User) *Profile {
    return &Profile{
        UserID: userId(u),
        Founds: StartingFounds,
    }
}

func GetProfile(req *http.Request) (profile *Profile, err error) {
    if val, ok := context.GetOk(req, reqProfile); ok {
        return val.(*Profile), nil
    }

    ctx := appengine.NewContext(req)

    u := user.Current(ctx)
    if u == nil {
        return nil, errors.New("User not logged in")
    }

    key := profileKeyForUser(ctx, u)

    profile = new(Profile)
    if datastore.Get(ctx, key, profile) == datastore.ErrNoSuchEntity {
        profile = newUserProfile(u)
        err = profile.Put(req)
    }

    context.Set(req, reqProfile, profile)
    return
}

func getProfileForUserId(req *http.Request, userID string) (profile *Profile, err error) {
    ctx := appengine.NewContext(req)
    key := profileKey(ctx, userID)

    profile = new(Profile)
    err = datastore.Get(ctx, key, profile)
    return
}
