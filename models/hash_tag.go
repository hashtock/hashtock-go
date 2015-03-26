package models

import (
    "gopkg.in/mgo.v2/bson"
)

const (
    HashTagCollectionName = "HashTag"
    initialInBankValue    = 100.0
)

type HashTag struct {
    HashTag string  `bson:"hashtag,omitempty" json:"hashtag"`
    Value   float64 `bson:"value,omitempty" json:"value"`
    InBank  float64 `bson:"in_bank,omitempty" json:"in_bank"`
}

func hashTagUpdateInBank(hashTag *HashTag, delta float64) (err error) {
    col := storage.Collection(HashTagCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "hashtag": hashTag.HashTag,
    }

    update_with := bson.M{
        "$inc": bson.M{"in_bank": delta},
    }

    err = col.Update(selector, update_with)
    return
}
