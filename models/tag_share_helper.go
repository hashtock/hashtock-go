package models

import (
    "net/http"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    "github.com/hashtock/hashtock-go/core"
)

func GetProfileShares(req *http.Request, profile *Profile) (shares []TagShare, err error) {
    col := storage.Collection(TagShareCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{"user_id": profile.UserID}
    err = col.Find(selector).Sort("hashtag").All(&shares)
    return
}

func GetProfileShareByTagName(req *http.Request, profile *Profile, hashTagName string) (tagShare *TagShare, err error) {
    tagShare, err = getOrCreateTagShare(req, profile, hashTagName)

    if (err == nil && tagShare.Quantity <= 0) || err == mgo.ErrNotFound {
        tagShare = nil
        err = core.NewNotFoundError(http.StatusText(http.StatusNotFound))
    }
    return
}

func getOrCreateTagShare(req *http.Request, profile *Profile, hashTagName string) (tagShare *TagShare, err error) {
    col := storage.Collection(TagShareCollectionName)
    defer col.Database.Session.Close()

    tagShare = &TagShare{
        HashTag: hashTagName,
        UserID:  profile.UserID,
    }

    selector := bson.M{
        "hashtag": hashTagName,
        "user_id": profile.UserID,
    }

    err = col.Find(selector).One(&tagShare)
    return
}
