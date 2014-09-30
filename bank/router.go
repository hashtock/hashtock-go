package bank

import (
    "github.com/gorilla/mux"
)

func AttachViews(r *mux.Router) {
    r.HandleFunc("/info/", HashTagInfoAllView)
    r.HandleFunc("/info/{hash}/", HashTagInfoView)
    r.HandleFunc("/sell_to_user/{hash}/{amount}/", SellToUserView)
    r.HandleFunc("/buy_from_user/{hash}/{amount}/", BuyFromUserView)
}
