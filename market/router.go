package market

import (
    "github.com/gorilla/mux"
)

func AttachViews(r *mux.Router) {
    r.HandleFunc("/", ViewOrders).Methods("GET")
    r.HandleFunc("/", AcceptNewOrder).Methods("POST", "PUT")
    r.HandleFunc("/{uuid}/", ViewOrder).Methods("GET")
    r.HandleFunc("/{uuid}/", DeleteOrder).Methods("DELETE")
}
