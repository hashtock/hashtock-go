package models

import (
    "gopkg.in/mgo.v2/bson"
)

const (
    TagShareCollectionName = "TagShare"
    minShareStep           = 0.01
)

type TagShare struct {
    HashTag  string  `bson:"hashtag" json:"hashtag"`
    UserID   string  `bson:"user_id" json:"-"`
    Quantity float64 `bson:"quantity" json:"quantity"`
}

func tagShareUpdateQuantity(tagShare *TagShare, delta float64) (err error) {
    col := storage.Collection(TagShareCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "hashtag": tagShare.HashTag,
        "user_id": tagShare.UserID,
    }

    update_with := bson.M{
        "$inc": bson.M{"quantity": delta},
    }
    // Upsert because it's possible that it does not exist
    _, err = col.Upsert(selector, update_with)
    return
}

func tagShareDelete(tagShare *TagShare) (err error) {
    col := storage.Collection(TagShareCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "hashtag": tagShare.HashTag,
        "user_id": tagShare.UserID,
    }

    err = col.Remove(selector)
    return
}
