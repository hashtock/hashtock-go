package models

import (
    "fmt"
    "net/http"
    "strings"

    "appengine"
    "appengine/datastore"

    "github.com/hashtock/hashtock-go/http_utils"
)

// User part of Order
type OrderBase struct {
    Action    string  `json:"action"`
    BankOrder bool    `json:"bank_order"`
    HashTag   string  `json:"hashtag"`
    Quantity  float64 `json:"quantity"`
}

// System fields regarding Order
// Read only for users
type OrderSystem struct {
    UUID     string `json:"uuid"`
    UserID   string `json:"user_id"`
    Complete bool   `json:"complete"`
}

type Order struct {
    OrderBase
    OrderSystem
}

const (
    orderKind  = "Order"
    actionBuy  = "buy"
    actionSell = "sell"
)

func (o *Order) key(ctx appengine.Context) (key *datastore.Key) {
    return orderKey(ctx, o.UUID)
}

func (o *Order) Put(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := o.key(ctx)
    _, err = datastore.Put(ctx, key, o)
    return
}

func (o *OrderBase) IsValid(req *http.Request) (err error) {
    fields := []string{}

    if (o.Action != actionBuy) && (o.Action != actionSell) {
        fields = append(fields, "action")
    }

    if !hashTagExists(req, o.HashTag) {
        fields = append(fields, "hashtag")
    }

    if o.Quantity <= 0 || o.Quantity > 100 {
        fields = append(fields, "quantity")
    }

    if len(fields) > 0 {
        msg := fmt.Sprintf("Incorrect fields: %s", strings.Join(fields, ", "))
        err = http_utils.NewBadRequestError(msg)
    }
    return
}
