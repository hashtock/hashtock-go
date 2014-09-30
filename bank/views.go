package bank

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "appengine"
    "github.com/gorilla/mux"
)

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

func SellToUserView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    vars := mux.Vars(req)
    hash := vars["hash"]
    amount, _ := strconv.ParseFloat(vars["amount"], 64)

    if err := SellToUser(ctx, hash, amount); err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    rw.WriteHeader(http.StatusOK)
    fmt.Fprintf(rw, "Sell: %v x %v\n", hash, amount)
}

func BuyFromUserView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    vars := mux.Vars(req)
    hash := vars["hash"]
    amount, _ := strconv.ParseFloat(vars["amount"], 64)

    if err := BuyFromUser(ctx, hash, amount); err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

    rw.WriteHeader(http.StatusOK)
    fmt.Fprintf(rw, "Sell: %v x %v\n", hash, amount)
}
