package models

import (
    "errors"

    "gopkg.in/mgo.v2"
)

type MgoStorage struct {
    session *mgo.Session
    db      string
    dbName  string
}

var storage *MgoStorage

func InitMongoStorage(dbUrl string, dbName string) (*MgoStorage, error) {
    if dbUrl == "" {
        return nil, errors.New("Url to Mongodb not provided")
    } else if dbName == "" {
        return nil, errors.New("Name of database for Mongodb not provided")
    }

    msession, err := mgo.Dial(dbUrl)
    if err != nil {
        return nil, err
    }

    storage = &MgoStorage{
        db:      dbUrl,
        dbName:  dbName,
        session: msession,
    }

    return storage, nil
}

func (m *MgoStorage) Collection(collectionName string) *mgo.Collection {
    session := m.session.Copy()
    return session.DB(m.dbName).C(collectionName)
}
