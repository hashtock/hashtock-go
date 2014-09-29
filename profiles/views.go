package profiles

import (
    "encoding/json"
    "net/http"

    "appengine"
    "github.com/gorilla/mux"
)

func AttachViews(r *mux.Router) {
    r.HandleFunc("/user/shares/", UserSharesView)
}

func UserSharesView(rw http.ResponseWriter, req *http.Request) {
    ctx := appengine.NewContext(req)

    profile, err := CurrentProfile(ctx)

    shares, err := profile.Shares(ctx)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    data, err := json.Marshal(shares)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(http.StatusOK)
    rw.Write(data)
}
