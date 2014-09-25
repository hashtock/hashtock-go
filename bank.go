package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "appengine"
    "appengine/datastore"
    "appengine/user"
    "github.com/gorilla/mux"
)

type BankEntry struct {
    Hash   string
    Value  float64
    InBank float64
}

// Returns info about given Tag:
// - it's bank value
// - how much is in Bank
func HashTagInfo(ctx appengine.Context, hash string) (entry BankEntry, err error) {
    key := datastore.NewKey(ctx, "BankEntry", hash, 0, nil)
    err = datastore.Get(ctx, key, &entry)
    return
}

// Returns info about all known Tags:
// - it's bank value
// - how much is in Bank
func HashTagInfoAll(ctx appengine.Context) (entries []BankEntry, err error) {
    q := datastore.NewQuery("BankEntry").Order("-Value")
    _, err = q.GetAll(ctx, &entries)

    return
}

func AddHashTag(ctx appengine.Context, hash string) (err error) {
    key := datastore.NewKey(ctx, "BankEntry", hash, 0, nil)

    entry := BankEntry{
        Hash:   hash,
        Value:  1,
        InBank: 100,
    }

    if datastore.Get(ctx, key, &entry) != datastore.ErrNoSuchEntity {
        err = errors.New("Tag already exist")
        return
    }

    _, err = datastore.Put(ctx, key, &entry)
    return
}

/*###########################################################*/

func HashTagInfoView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    vars := mux.Vars(req)

    hash := vars["hash"]

    entry, err := HashTagInfo(ctx, hash)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
        return
    }

    data, err := json.Marshal(entry)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(http.StatusOK)
    fmt.Fprintln(rw, string(data))
}

func HashTagInfoAllView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    entries, err := HashTagInfoAll(ctx)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
        return
    }

    data, err := json.Marshal(entries)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(http.StatusOK)
    fmt.Fprintln(rw, string(data))
}

func AddHashTagView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)
    if !user.IsAdmin(ctx) {
        http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
        return
    }

    vars := mux.Vars(req)

    hash := vars["hash"]
    if err := AddHashTag(ctx, hash); err != nil {
        http.Error(rw, err.Error(), http.StatusNotFound)
    }

    rw.WriteHeader(http.StatusOK)
    fmt.Fprintf(rw, "OK: %#v", hash)
}

///////////////////////////////////////////

func SellToBank(hash string, amount float64) (err error) {
    if 0 >= amount || amount > 100 {
        err = errors.New("Amount outside boundries")
        return
    }

    return
}

func SellView(rw http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)

    amount, _ := strconv.ParseFloat(vars["amount"], 64)
    hash := vars["hash"]

    err := SellToBank(hash, amount)

    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Fprintf(rw, "Sell %v\n", vars["hash"])
}
