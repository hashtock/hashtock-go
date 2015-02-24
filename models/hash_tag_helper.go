package models

import (
    "fmt"
    "net/http"
    "strings"

    "appengine"
    "appengine/datastore"
    "appengine/user"

    "github.com/hashtock/hashtock-go/core"
)

func hashTagKey(ctx appengine.Context, hash_tag string) (key *datastore.Key) {
    return datastore.NewKey(ctx, hashTagKind, hash_tag, 0, nil)
}

func GetAllHashTags(req *http.Request) (hash_tags []*HashTag, err error) {
    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(hashTagKind).Order("-Value")
    _, err = q.GetAll(ctx, &hash_tags)
    return
}

func GetHashTag(req *http.Request, hash_tag_name string) (hash_tag *HashTag, err error) {
    ctx := appengine.NewContext(req)

    hash_tag = new(HashTag)
    key := hashTagKey(ctx, hash_tag_name)
    err = datastore.Get(ctx, key, hash_tag)

    if err == datastore.ErrNoSuchEntity {
        msg := fmt.Sprintf("HashTag %#v not found", hash_tag_name)
        err = core.NewNotFoundError(msg)
    } else if err != nil {
        err = core.NewInternalServerError(err.Error())
    }

    return
}

func GetOrCreateHashTag(req *http.Request, hashTagName string) (hashTag *HashTag, err error) {
    if err = CanUserCreateUpdateHashTag(req); err != nil {
        return
    }

    ctx := appengine.NewContext(req)

    hashTag = &HashTag{
        HashTag: hashTagName,
    }

    key := hashTagKey(ctx, hashTagName)
    err = datastore.Get(ctx, key, hashTag)

    if err == datastore.ErrNoSuchEntity {
        hashTag.InBank = initialInBankValue
        if err = hashTag.Put(req); err != nil {
            return
        }
    } else if err != nil {
        err = core.NewInternalServerError(err.Error())
    }

    return
}

func hashTagExists(req *http.Request, hash_tag_name string) (ok bool, err error) {
    if strings.TrimSpace(hash_tag_name) != hash_tag_name || (hash_tag_name == "") {
        return false, core.NewBadRequestError("Tag name invalid")
    }

    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(hashTagKind).Filter("HashTag =", hash_tag_name)

    var count int
    count, err = q.Count(ctx)

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

func CanUserCreateUpdateHashTag(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    if req.Header.Get("X-AppEngine-Cron") == "true" {
        return nil
    }

    if !user.IsAdmin(ctx) {
        return core.NewForbiddenError()
    }

    return nil
}

func AddHashTag(req *http.Request, tag HashTag) (new_tag HashTag, err error) {
    if err = CanUserCreateUpdateHashTag(req); err != nil {
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
    if err = tag.Put(req); err != nil {
        return
    }

    return tag, err
}

func UpdateHashTagValue(req *http.Request, hash_tag_name string, new_value float64) (tag *HashTag, err error) {
    if err = CanUserCreateUpdateHashTag(req); err != nil {
        return
    }

    if new_value <= 0 {
        err = core.NewBadRequestError("Value has to be positive")
        return
    }

    tag, err = GetHashTag(req, hash_tag_name)
    if err != nil {
        return
    }

    if tag.Value == new_value {
        return
    }

    tag.Value = new_value
    err = tag.Put(req)
    return
}
