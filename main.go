package main

import (
	"log"
	"net/http"

	authClient "github.com/hashtock/auth/client"
	"github.com/hashtock/service-tools/serialize"
	"github.com/nats-io/nats"

	"github.com/hashtock/hashtock-go/conf"
	"github.com/hashtock/hashtock-go/jobs"
	"github.com/hashtock/hashtock-go/models"
	"github.com/hashtock/hashtock-go/webapp"
)

func main() {
	cfg := conf.GetConfig()
	storage, err := models.InitMongoStorage(cfg.General.DB, cfg.General.DBName)
	if err != nil {
		log.Fatalln(err)
	}

	tagWorker := jobs.NewTagValueWorker(storage, nats.DefaultURL, "tags.counts")
	tagWorker.Start()
	defer tagWorker.Stop()

	orderWorker := jobs.NewOrderWorker(storage, storage, storage, cfg.Jobs.BankOrders)
	orderWorker.Start()
	defer orderWorker.Stop()

	whoClient, whoErr := authClient.NewClient(cfg.General.AuthAddress)
	if whoErr != nil {
		log.Fatalln(whoErr)
	}
	options := webapp.Options{
		Serializer:       &serialize.WebAPISerializer{},
		PortfolioStorage: storage,
		BankStorage:      storage,
		OrderStorage:     storage,
		WhoClient:        whoClient,
	}

	handler := webapp.Handlers(options)

	err = http.ListenAndServe(cfg.General.ServeAddr, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
