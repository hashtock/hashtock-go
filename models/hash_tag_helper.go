package models

import (
    "fmt"
    "net/http"
    "strings"

    "appengine"
    "appengine/datastore"
    "appengine/user"

    "github.com/hashtock/hashtock-go/http_utils"
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
        err = http_utils.NewNotFoundError(msg)
    } else if err != nil {
        err = http_utils.NewInternalServerError(err.Error())
    }

    return
}

func hashTagExists(req *http.Request, hash_tag_name string) (ok bool, err error) {
    if strings.TrimSpace(hash_tag_name) != hash_tag_name || (hash_tag_name == "") {
        return false, http_utils.NewBadRequestError("Tag name invalid")
    }

    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(hashTagKind).Filter("HashTag =", hash_tag_name)

    var count int
    count, err = q.Count(ctx)

    return count > 0 || err != nil, err
}

func CanUserCreateHashTag(req *http.Request) (err error) {
    ctx := appengine.NewContext(req)

    if !user.IsAdmin(ctx) {
        msg_404 := http.StatusText(http.StatusNotFound)
        return http_utils.NewNotFoundError(msg_404)
    }

    return nil
}

func AddHashTag(req *http.Request, tag HashTag) (new_tag HashTag, err error) {
    if err = CanUserCreateHashTag(req); err != nil {
        return
    }

    var exists bool
    if exists, err = hashTagExists(req, tag.HashTag); err != nil {
        return
    } else if exists {
        err = http_utils.NewBadRequestError("Tag alread exists")
        return
    }

    tag.InBank = 100.0
    if err = tag.Put(req); err != nil {
        return
    }

    return tag, err
}
