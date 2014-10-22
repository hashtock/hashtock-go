package models

import (
    "fmt"
    "net/http"
    "strings"

    "appengine"
    "appengine/datastore"

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

func hashTagExists(req *http.Request, hash_tag_name string) bool {
    if strings.TrimSpace(hash_tag_name) != hash_tag_name || (hash_tag_name == "") {
        return false
    }

    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(hashTagKind).Filter("HashTag =", hash_tag_name)

    count, err := q.Count(ctx)

    return count > 0 || err != nil
}
