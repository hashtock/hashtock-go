package main

import (
    "log"
    "net/http"

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

    jobs.StartJobs(cfg, storage)

    handler := webapp.Handlers(cfg, storage)

    err = http.ListenAndServe(cfg.General.ServeAddr, handler)
    if err != nil {
        log.Fatalln(err)
    }
}
