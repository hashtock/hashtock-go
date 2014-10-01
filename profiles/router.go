package profiles

import (
    "github.com/gorilla/mux"
)

func AttachViews(r *mux.Router) {
    r.HandleFunc("/user/shares/", UserSharesView)
}
