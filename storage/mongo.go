package storage

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/hashtock-go/core"
)

const (
	BankCollectionName  = "Bank"
	OrderCollectionName = "Orders"

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

func (m *MgoStorage) collection(collectionName string) *mgo.Collection {
	lsession := m.session.Copy()
	col := lsession.DB(m.dbName).C(collectionName)
	return col
}

// Helpers

func (m *MgoStorage) inBank(hashTags ...string) (tagsInBank map[string]float64, err error) {
	var results = []struct {
		OnMarket float64 `bson:"on_market"`
		HashTag  string  `bson:"hashtag"`
	}{}

	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"hashtag": bson.M{
					"$in": hashTags,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$hashtag",
				"on_market": bson.M{
					"$sum": "$quantity",
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":       0,
				"hashtag":   "$_id",
				"on_market": 1,
			},
		},
	}

	pipe := col.Pipe(pipeline)
	err = pipe.All(&results)
	if err != nil {
		return
	}

	tagsInBank = make(map[string]float64, len(hashTags))
	for _, tag := range hashTags {
		tagsInBank[tag] = initialInBankValue
	}

	for _, result := range results {
		tagsInBank[result.HashTag] -= result.OnMarket
	}

	return
}

func (m *MgoStorage) portfolioPipeline(userId string, hashTag string, byTag bool) []bson.M {
	selector := bson.M{
		"user_id":    userId,
		"resolution": core.SUCCESS,
	}

	groupBy := bson.M{
		"user_id": "$user_id",
	}

	project := bson.M{
		"_id":      0,
		"user_id":  "$_id.user_id",
		"quantity": 1,
		"value":    1,
	}

	if hashTag != "" {
		selector["hashtag"] = hashTag
	}

	if byTag {
		groupBy["hashtag"] = "$hashtag"
		project["hashtag"] = "$_id.hashtag"
	}

	pipeline := []bson.M{
		bson.M{
			"$match": selector,
		},

		bson.M{
			"$group": bson.M{
				"_id": groupBy,
				"quantity": bson.M{
					"$sum": "$quantity",
				},
			},
		},

		bson.M{
			"$project": project,
		},

		bson.M{
			"$sort": bson.M{"hashtag": 1},
		},
	}

	return pipeline
}

// Portfolio interface

func (m *MgoStorage) PortfolioShares(userId string) (shares []core.TagShare, err error) {
	col := m.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	pipeline := m.portfolioPipeline(userId, "", true)
	pipe := col.Pipe(pipeline)
	err = pipe.All(&shares)
	return
}

func (m *MgoStorage) PortfolioShare(userId string, hashTagName string, strict bool) (tagShare *core.TagShare, err error) {
	col := storage.collection(OrderCollectionName)
	defer col.Database.Session.Close()

	pipeline := m.portfolioPipeline(userId, hashTagName, true)
	pipe := col.Pipe(pipeline)
	err = pipe.One(&tagShare)

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

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"user_id":    userId,
				"resolution": core.SUCCESS,
			},
		},

		bson.M{
			"$group": bson.M{
				"_id": "$user_id",
				"cash": bson.M{
					"$sum": "$value",
				},
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

// Bank interface

func (m *MgoStorage) Tags() (hashTags []core.HashTag, err error) {
	col := m.collection(BankCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(nil).Sort("-value").All(&hashTags)
	if err != nil {
		return
	}

	var tags []string = make([]string, len(hashTags))

	for i, tag := range hashTags {
		tags[i] = tag.HashTag
	}

	tagsInBank, err := m.inBank(tags...)
	if err != nil {
		return
	}

	for i, tag := range hashTags {
		hashTags[i].InBank = tagsInBank[tag.HashTag]
	}

	return
}

func (m *MgoStorage) Tag(hashTagName string) (hashTag *core.HashTag, err error) {
	col := m.collection(BankCollectionName)
	defer col.Database.Session.Close()

	selector := bson.M{
		"hashtag": hashTagName,
	}

	err = col.Find(selector).One(&hashTag)

	if err == mgo.ErrNotFound {
		msg := fmt.Sprintf("HashTag %#v not found", hashTagName)
		err = core.NewNotFoundError(msg)
	} else if err != nil {
		err = core.NewInternalServerError(err.Error())
	}

	if err != nil {
		return
	}

	tagsInBank, err := m.inBank(hashTagName)
	if err != nil {
		return
	}

	hashTag.InBank = tagsInBank[hashTagName]

	return
}

func (m *MgoStorage) TagSetValue(hashTagName string, value float64) error {
	var hashTag core.HashTag
	if value < 0 {
		return core.NewBadRequestError("Value can not be negative")
	}

	col := m.collection(BankCollectionName)
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
		}
		return col.Insert(hashTag)
	}

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
