package models

import (
    "gopkg.in/mgo.v2/bson"
)

const (
    ProfileCollectionName = "Profile"
    StartingFounds        = 1000
)

type Profile struct {
    UserID  string  `bson:"user_id,omitempty" json:"id"`
    Founds  float64 `bson:"founds,omitempty" json:"founds"`
    IsAdmin bool    `bson:"is_admin,omitempty" json:"-"`
}

func profileUpdateFounds(profile *Profile, delta float64) (err error) {
    col := storage.Collection(ProfileCollectionName)
    defer col.Database.Session.Close()

    selector := bson.M{
        "user_id": profile.UserID,
    }

    update_with := bson.M{
        "$inc": bson.M{"founds": delta},
    }
    err = col.Update(selector, update_with)
    return
}
