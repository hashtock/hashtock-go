package models

import (
    "fmt"
    "net/http"

    "appengine"
    "appengine/datastore"

    "github.com/hashtock/hashtock-go/core"
)

func tagShareKey(ctx appengine.Context, hashTagName, userId string) (key *datastore.Key) {
    share_id := fmt.Sprintf("%s-%s", hashTagName, userId)
    return datastore.NewKey(ctx, tagShareKind, share_id, 0, nil)
}

func GetProfileShares(req *http.Request, profile *Profile) (shares []TagShare, err error) {
    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(tagShareKind)
    q = q.Filter("UserID =", profile.UserID).Filter("Quantity >=", 0.0)
    q = q.Order("Quantity").Order("HashTag")
    _, err = q.GetAll(ctx, &shares)
    return
}

func GetProfileShareByTagName(req *http.Request, profile *Profile, hashTagName string) (tagShare *TagShare, err error) {
    tagShare, err = getOrCreateTagShare(req, profile, hashTagName)

    if err == nil && tagShare.Quantity <= 0 {
        tagShare = nil
        err = core.NewNotFoundError(http.StatusText(http.StatusNotFound))
    }
    return
}

func getOrCreateTagShare(req *http.Request, profile *Profile, hashTagName string) (tagShare *TagShare, err error) {
    ctx := appengine.NewContext(req)

    key := tagShareKey(ctx, hashTagName, profile.UserID)

    tagShare = new(TagShare)
    err = datastore.Get(ctx, key, tagShare)
    if err == datastore.ErrNoSuchEntity {
        err = nil
        tagShare = &TagShare{
            HashTag: hashTagName,
            UserID:  profile.UserID,
        }
    }

    return
}
