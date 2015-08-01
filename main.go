package main

import (
	"log"
	"net/http"

	authClient "github.com/hashtock/auth/client"
	"github.com/hashtock/service-tools/serialize"
	"github.com/nats-io/nats"

	"github.com/hashtock/hashtock-go/conf"
	"github.com/hashtock/hashtock-go/jobs"
	"github.com/hashtock/hashtock-go/storage"
	"github.com/hashtock/hashtock-go/webapp"
)

func main() {
	cfg := conf.GetConfig()
	dataStore, err := storage.InitMongoStorage(cfg.General.DB, cfg.General.DBName)
	if err != nil {
		log.Fatalln(err)
	}

	tagWorker := jobs.NewTagValueWorker(dataStore, nats.DefaultURL, "tags.counts")
	tagWorker.Start()
	defer tagWorker.Stop()

	orderWorker := jobs.NewOrderWorker(dataStore, dataStore, dataStore, cfg.Jobs.BankOrders)
	orderWorker.Start()
	defer orderWorker.Stop()

	whoClient, whoErr := authClient.NewClient(cfg.General.AuthAddress)
	if whoErr != nil {
		log.Fatalln(whoErr)
	}
	options := webapp.Options{
		Serializer:       &serialize.WebAPISerializer{},
		PortfolioStorage: dataStore,
		BankStorage:      dataStore,
		OrderStorage:     dataStore,
		WhoClient:        whoClient,
	}

	handler := webapp.Handlers(options)

	err = http.ListenAndServe(cfg.General.ServeAddr, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
