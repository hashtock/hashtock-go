package models

import (
    "time"
)

const (
    HashTagValueCollectionName = "HashTagValue"
)

type HashTagValue struct {
    HashTag string    `bson:"hashtag,omitempty" json:"-"`
    Value   float64   `bson:"value,omitempty" json:"value"`
    Date    time.Time `bson:"date,omitempty" json:"date"`
}

type byDate []HashTagValue

func (a byDate) Len() int           { return len(a) }
func (a byDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
