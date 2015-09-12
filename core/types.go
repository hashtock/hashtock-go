package core

import (
	"time"
)

type OrderResolution string
type OrderType string
type OrderAction string

const (
	PENDING OrderResolution = ""
	SUCCESS                 = "success"
	FAILURE                 = "failure"
	ERROR                   = "error"

	TYPE_BANK          OrderType = "bank"
	TYPE_MARKET        OrderType = "market"
	TYPE_MARKET_FULFIL OrderType = "market_fulfil"
)

type Balance struct {
	Cash float64 `bson:"cash" json:"cash"`
}

type HashTag struct {
	HashTag string  `bson:"hashtag,omitempty" json:"hashtag"`
	Value   float64 `bson:"value,omitempty" json:"value"`
	InBank  float64 `bson:"-" json:"in_bank"`
}

type TagShare struct {
	HashTag  string  `bson:"hashtag" json:"hashtag"`
	UserID   string  `bson:"user_id" json:"-"`
	Quantity float64 `bson:"quantity" json:"quantity"`
}

// User part of Order
type OrderBase struct {
	Type        OrderType `bson:"type" json:"type"`
	HashTag     string    `bson:"hashtag" json:"hashtag"`
	Quantity    float64   `bson:"quantity" json:"quantity"`
	UnitPrice   float64   `bson:"unit_price" json:"unit_price,omitempty"`
	BaseOrderID string    `bson:"base_order_id" json:"base_order_id,omitempty"`
}

// System fields regarding Order
// Read only for users
type OrderSystem struct {
	UUID       string          `bson:"uuid" json:"uuid"`
	UserID     string          `bson:"user_id" json:"-"`
	Complete   bool            `bson:"complete" json:"complete"`
	Value      float64         `bson:"value" json:"value"`
	CreatedAt  time.Time       `bson:"created_at" json:"created_at"`
	ExecutedAt time.Time       `bson:"executed_at" json:"executed_at"`
	Resolution OrderResolution `bson:"resolution" json:"resolution"`
	Notes      string          `bson:"notes" json:"notes"`
}

type Order struct {
	OrderBase   `bson:",inline"`
	OrderSystem `bson:",inline"`
}

type OrderFilter struct {
	UserID   string `bson:"user_id,omitempty" schema:"-"`
	Complete bool   `bson:"complete" schema:"-"`

	Type       OrderType       `bson:"type,omitempty" schema:"type"`
	HashTag    string          `bson:"hashtag,omitempty" schema:"tag"`
	Resolution OrderResolution `bson:"resolution,omitempty" schema:"resolution"`
}
