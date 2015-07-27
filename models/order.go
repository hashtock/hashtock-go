package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/hashtock-go/core"
)

// type OrderResolution string

// const (
//     PENDING OrderResolution = ""
//     SUCCESS                 = "success"
//     FAILURE                 = "failure"
//     ERROR                   = "error"
// )

// // User part of Order
// type OrderBase struct {
//     Action    string  `bson:"action" json:"action"`
//     BankOrder bool    `bson:"bank_order" json:"bank_order"`
//     HashTag   string  `bson:"hashtag" json:"hashtag"`
//     Quantity  float64 `bson:"quantity" json:"quantity"`
// }

// // System fields regarding Order
// // Read only for users
// type OrderSystem struct {
//     UUID       string          `bson:"uuid" json:"uuid"`
//     UserID     string          `bson:"user_id" json:"user_id"`
//     Complete   bool            `bson:"complete" json:"complete"`
//     Value      float64         `bson:"value" json:"value"`
//     CreatedAt  time.Time       `bson:"created_at" json:"created_at"`
//     ExecutedAt time.Time       `bson:"executed_at" json:"executed_at"`
//     Resolution OrderResolution `bson:"resolution" json:"resolution"`
//     Notes      string          `bson:"notes" json:"notes"`
// }

// type Order struct {
// 	OrderBase   `bson:",inline"`
// 	OrderSystem `bson:",inline"`
// }

const (
	OrderCollectionName = "Order"
	actionBuy           = "buy"
	actionSell          = "sell"
)

func orderDelete(order *core.Order) (err error) {
	col := storage.Collection(OrderCollectionName)
	defer col.Database.Session.Close()

	err = col.Remove(order)
	return
}

func baseOrderValid(req *http.Request, order core.OrderBase) (err error) {
	fields := []string{}

	// if (order.Action != actionBuy) && (order.Action != actionSell) {
	// 	fields = append(fields, "action")
	// }

	if exists, tmp_err := hashTagExists(req, order.HashTag); !exists || tmp_err != nil {
		fields = append(fields, "hashtag")
	}

	if order.Quantity < minShareStep || order.Quantity > 100 {
		fields = append(fields, "quantity")
	}

	if len(fields) > 0 {
		msg := fmt.Sprintf("Incorrect fields: %s", strings.Join(fields, ", "))
		err = core.NewBadRequestError(msg)
	}
	return
}

func markOrderAsComplete(order core.Order, status core.OrderResolution, notes string) (err error) {
	order.Complete = true
	order.Resolution = status
	order.Notes = notes
	order.ExecutedAt = time.Now()

	col := storage.Collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"uuid": order.UUID,
	}
	err = col.Update(selector, order)
	return
}
