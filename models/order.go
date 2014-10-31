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

func (o *Order) Delete(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    key := o.key(ctx)
    err = datastore.Delete(ctx, key)
    return
}

func (o *OrderBase) IsValid(req *http.Request) (err error) {
    fields := []string{}

    if (o.Action != actionBuy) && (o.Action != actionSell) {
        fields = append(fields, "action")
    }

    if exists, tmp_err := hashTagExists(req, o.HashTag); !exists || tmp_err != nil {
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

func (o *Order) canAccess(req *http.Request) (ok bool, err error) {
    var profile *Profile

    if profile, err = GetProfile(req); err != nil {
        return
    }

    ok = o.UserID == profile.UserID
    return
}

func (o *Order) isCancellable() bool {
    return o.Complete == false
}

func (o *Order) isBuy() bool {
    return o.Action == actionBuy
}

func (o *Order) markAsComplete(req *http.Request) (err error) {
    o.Complete = true
    return o.Put(req)
}
