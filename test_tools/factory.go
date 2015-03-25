package test_tools

import (
    "fmt"

    "github.com/hashtock/hashtock-go/models"
)

func (t *TestApp) CreateProfile(name string, is_admin bool) (profile *models.Profile, err error) {
    col := t.storage.Collection(models.ProfileCollectionName)
    defer col.Database.Session.Close()

    profile = &models.Profile{
        UserID:  name,
        IsAdmin: is_admin,
        Founds:  models.StartingFounds,
    }
    err = col.Insert(&profile)
    return
}

func (t *TestApp) put(obj interface{}) (err error) {
    colName := ""
    switch obj.(type) {
    case models.HashTag:
        colName = models.HashTagCollectionName
    case models.HashTagValue:
        colName = models.HashTagValueCollectionName
    case models.Order:
        colName = models.OrderCollectionName
    case models.TagShare:
        colName = models.TagShareCollectionName
    default:
        msg := fmt.Sprintf("Type %T not supported", obj)
        panic(msg)
    }

    col := t.storage.Collection(colName)
    defer col.Database.Session.Close()
    return col.Insert(obj)
}

func (t *TestApp) Put(obj ...interface{}) (err error) {
    for _, o := range obj {
        err = t.put(o)
        if err != nil {
            return
        }
    }
    return
}
