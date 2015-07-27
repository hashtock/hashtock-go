package models

/*
ToDo:
- Refactor Get, Create, GetOrCreate to reduce amount of code
*/

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/hashtock/hashtock-go/core"
)

func GetAllHashTags(req *http.Request) (hashTags []core.HashTag, err error) {
	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(nil).Sort("-value").All(&hashTags)
	return
}

func GetHashTag(req *http.Request, hash_tag_name string) (hash_tag *core.HashTag, err error) {
	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(bson.M{"hashtag": hash_tag_name}).One(&hash_tag)

	if err == mgo.ErrNotFound {
		msg := fmt.Sprintf("HashTag %#v not found", hash_tag_name)
		err = core.NewNotFoundError(msg)
	} else if err != nil {
		err = core.NewInternalServerError(err.Error())
	}

	return
}

//ToDo: Needs tests
func GetOrCreateHashTag(req *http.Request, profile *Profile, hashTagName string) (hashTag *core.HashTag, err error) {
	if err = CanUserCreateUpdateHashTag(req, profile); err != nil {
		return
	}

	hashTag = &core.HashTag{
		HashTag: hashTagName,
	}

	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	err = col.Find(hashTag).One(&hashTag)

	if err == mgo.ErrNotFound {
		hashTag.InBank = initialInBankValue
		_, err = col.Upsert(hashTag, hashTag)
	}

	if err != nil {
		err = core.NewInternalServerError(err.Error())
	}

	return
}

func hashTagExists(req *http.Request, hashTagName string) (ok bool, err error) {
	if strings.TrimSpace(hashTagName) != hashTagName || (hashTagName == "") {
		return false, core.NewBadRequestError("Tag name invalid")
	}

	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	selector := core.HashTag{HashTag: hashTagName}

	var count int
	count, err = col.Find(selector).Count()

	return count > 0 || err != nil, err
}

func hashTagExistsOrError(req *http.Request, hash_tag_name string) (err error) {
	var exists bool

	exists, err = hashTagExists(req, hash_tag_name)
	if err != nil {
		return
	}

	if !exists {
		msg := fmt.Sprintf("Tag '%v' does not exist", hash_tag_name)
		return core.NewNotFoundError(msg)
	}
	return
}

func CanUserCreateUpdateHashTag(req *http.Request, profile *Profile) (err error) {
	if profile != nil && !profile.IsAdmin {
		return core.NewForbiddenError()
	}

	return
}

func AddHashTag(req *http.Request, profile *Profile, tag core.HashTag) (newTag core.HashTag, err error) {
	if err = CanUserCreateUpdateHashTag(req, profile); err != nil {
		return
	}

	var exists bool
	if exists, err = hashTagExists(req, tag.HashTag); err != nil {
		return
	} else if exists {
		err = core.NewBadRequestError("Tag alread exists")
		return
	}

	tag.InBank = initialInBankValue

	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	if _, err := col.Upsert(tag, tag); err != nil {
		return tag, err
	}

	return tag, err
}

func UpdateHashTagValue(req *http.Request, profile *Profile, hash_tag_name string, new_value float64) (tag core.HashTag, err error) {
	if err = CanUserCreateUpdateHashTag(req, profile); err != nil {
		return
	}

	if new_value <= 0 {
		err = core.NewBadRequestError("Value has to be positive")
		return
	}

	col := storage.Collection(HashTagCollectionName)
	defer col.Database.Session.Close()

	selector := core.HashTag{HashTag: hash_tag_name}
	change := mgo.Change{
		Update: bson.M{
			"$set": core.HashTag{Value: new_value},
		},
		ReturnNew: true,
	}
	_, err = col.Find(selector).Apply(change, &tag)

	if err == mgo.ErrNotFound {
		msg := fmt.Sprintf("HashTag %#v not found", hash_tag_name)
		err = core.NewNotFoundError(msg)
	}

	return
}
