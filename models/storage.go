package models

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/hashtock-go/core"
)

const (
	HashTagCollectionName  = "HashTag"
	OrderCollectionName    = "Order"
	TagShareCollectionName = "TagShare"

	initialInBankValue = 100.0
	StartingFounds     = 1000
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

// ToDo: Remove
func (m *MgoStorage) Collection(collectionName string) *mgo.Collection {
	session := m.session.Copy()
	return session.DB(m.dbName).C(collectionName)
}

func (m *MgoStorage) collection(collectionName string) *mgo.Collection {
	lsession := m.session.Copy()
	col := lsession.DB(m.dbName).C(collectionName)
	return col
}

// func (m *MgoStorage) getOrCreateTagShare(userId string, hashTagName string) (tagShare *core.TagShare, err error) {
// 	col := storage.Collection(TagShareCollectionName)
// 	defer col.Database.Session.Close()

// 	tagShare = &core.TagShare{
// 		HashTag: hashTagName,
// 		UserID:  userId,
// 	}

// 	selector := bson.M{
// 		"hashtag": hashTagName,
// 		"user_id": userId,
// 	}

// 	err = col.Find(selector).One(&tagShare)
// 	return
// }

// Portfolio interface

func (m *MgoStorage) PortfolioShares(userId string) (shares []core.TagShare, err error) {
	col := m.collection(TagShareCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{"user_id": userId}
	err = col.Find(selector).Sort("hashtag").All(&shares)
	return
}

func (m *MgoStorage) PortfolioShare(userId string, hashTagName string, strict bool) (tagShare *core.TagShare, err error) {
	col := storage.Collection(TagShareCollectionName)
	defer col.Database.Session.Close()

	tagShare = &core.TagShare{
		HashTag: hashTagName,
		UserID:  userId,
	}

	selector := bson.M{
		"hashtag": hashTagName,
		"user_id": userId,
	}

	err = col.Find(selector).One(&tagShare)
	if (err == nil && tagShare.Quantity <= 0) || err == mgo.ErrNotFound {
		if !strict {
			// Return empty share if not strict
			err = nil
		} else {
			tagShare = nil
			msg := fmt.Sprintf("User %#v does not own %#v tag shares", userId, hashTagName)
			err = core.NewNotFoundError(msg)
		}
	}

	return
}

func (m *MgoStorage) PortfolioBalance(userId string) (balance core.Balance, err error) {
	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	query := bson.M{
		"user_id":    userId,
		"resolution": core.SUCCESS,
	}

	pipeline := []bson.M{
		bson.M{"$match": query},

		bson.M{
			"$group": bson.M{
				"_id":  "$user_id",
				"cash": bson.M{"$sum": "$value"},
			},
		},

		bson.M{
			"$project": bson.M{
				"_id":  0,
				"cash": 1,
			},
		},
	}

	pipe := col.Pipe(pipeline)
	err = pipe.One(&balance)

	// No shares yet
	if err == mgo.ErrNotFound {
		err = nil
		balance.Cash = 0
	}

	balance.Cash += StartingFounds
	return
}

func (m *MgoStorage) PortfolioShareUpdateQuantity(userId string, tag string, quantity float64) error {
	var tagShare core.TagShare
	col := storage.Collection(TagShareCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"hashtag": tag,
		"user_id": userId,
	}

	query := col.Find(selector)
	cnt, err := query.Count()
	if err != nil {
		return err
	}

	// Need new tag
	if cnt == 0 {
		if quantity < 0 {
			return errors.New("Selling short is not allowed")
		}
		tagShare = core.TagShare{
			HashTag:  tag,
			UserID:   userId,
			Quantity: quantity,
		}
		return col.Insert(&tagShare)
	}

	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{"quantity": quantity},
		},
	}
	_, err = query.Apply(change, &tagShare)
	return err
}

// Bank interface

func (m *MgoStorage) Tags() (hashTags []core.HashTag, err error) {
	col := m.collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(nil).Sort("-value").All(&hashTags)
	return
}

func (m *MgoStorage) Tag(hashTagName string) (hashTag *core.HashTag, err error) {
	col := m.collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(bson.M{"hashtag": hashTagName}).One(&hashTag)

	if err == mgo.ErrNotFound {
		msg := fmt.Sprintf("HashTag %#v not found", hashTagName)
		err = core.NewNotFoundError(msg)
	} else if err != nil {
		err = core.NewInternalServerError(err.Error())
	}

	return
}

func (m *MgoStorage) TagSetValue(hashTagName string, value float64) error {
	var hashTag core.HashTag
	if value < 0 {
		return core.NewBadRequestError("Value can not be negative")
	}

	col := m.collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{"hashtag": hashTagName}
	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{"value": value},
		},
	}
	_, err := col.Find(selector).Apply(change, &hashTag)
	if err == mgo.ErrNotFound {
		hashTag = core.HashTag{
			HashTag: hashTagName,
			Value:   value,
			InBank:  initialInBankValue,
		}
		return col.Insert(hashTag)
	}

	return err
}

func (m *MgoStorage) TagUpdateInBank(tag string, quantity float64) error {
	var hashTag core.HashTag

	col := m.collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{"hashtag": tag}
	if err := col.Find(selector).One(&hashTag); err != nil {
		return err
	}

	newInBank := hashTag.InBank + quantity
	if newInBank < 0.0 {
		return errors.New("Not enough shares of in bank.")
	} else if newInBank > 100.0 {
		return errors.New("Bank can not own more then 100% of shres")
	}

	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{"in_bank": quantity},
		},
	}
	_, err := col.Find(selector).Apply(change, &hashTag)
	return err
}

// Order interface

func (m *MgoStorage) Order(userId string, orderId string) (order *core.Order, err error) {
	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"user_id": userId,
		"uuid":    orderId,
	}

	err = col.Find(selector).One(&order)
	if err == mgo.ErrNotFound {
		notFoundMsg := fmt.Sprintf("Order %#v not found", orderId)
		err = core.NewNotFoundError(notFoundMsg)
	}
	return
}

func (m *MgoStorage) Orders(userId string, complete bool, tag string, resolution string) (orders []core.Order, err error) {
	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"complete": complete,
		"user_id":  userId,
	}

	if tag != "" {
		selector["hashtag"] = tag
	}

	if resolution != "" {
		selector["resolution"] = resolution
	}

	err = col.Find(selector).Sort("-created_at").All(&orders)
	return
}

func (m *MgoStorage) AddOrder(order *core.Order) (err error) {
	col := storage.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	err = col.Insert(order)
	return
}

func (m *MgoStorage) DeleteOrder(userId string, orderId string) (err error) {
	col := storage.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"user_id": userId,
		"uuid":    orderId,
	}

	err = col.Remove(selector)
	return
}

// OrderExecuter interface

func (m *MgoStorage) OrdersToExecute() (orders []core.Order, err error) {
	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"complete": false,
	}

	err = col.Find(selector).Sort("-created_at").All(&orders)
	return
}

func (m *MgoStorage) OrderCompleted(orderId string, status core.OrderResolution, notes string) (err error) {
	var order core.Order
	col := storage.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"uuid": orderId,
	}
	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"complete":    true,
				"resolution":  status,
				"notes":       notes,
				"executed_at": time.Now(),
			},
		},
	}
	_, err = col.Find(selector).Apply(change, &order)
	return
}
