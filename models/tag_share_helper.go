package models

import (
    "net/http"

    "appengine"
    "appengine/datastore"
)

func tagShareKey(ctx appengine.Context, share_id string) (key *datastore.Key) {
    return datastore.NewKey(ctx, tagShareKind, share_id, 0, nil)
}

func GetProfileShares(req *http.Request, profile *Profile) (shares []TagShare, err error) {
    ctx := appengine.NewContext(req)

    q := datastore.NewQuery(tagShareKind).Filter("UserID =", profile.UserID).Order("HashTag")
    _, err = q.GetAll(ctx, &shares)
    return
}
