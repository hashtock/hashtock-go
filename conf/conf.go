package conf

import (
    "net/http"

    "appengine"
    "appengine/datastore"
)

const (
    confingKind = "Config"
)

type Config struct {
    TrackerHMACSecret string
    TrackerHost       string
}

func (c *Config) key(ctx appengine.Context) (key *datastore.Key) {
    return datastore.NewKey(ctx, confingKind, "Config", 0, nil)
}

func (c *Config) Put(req *http.Request) error {
    ctx := appengine.NewContext(req)
    _, err := datastore.Put(ctx, c.key(ctx), c)
    return err
}

func (c *Config) Get(req *http.Request) error {
    ctx := appengine.NewContext(req)
    err := datastore.Get(ctx, c.key(ctx), c)

    if err == datastore.ErrNoSuchEntity {
        return c.Put(req)
    }

    return err
}

func TrackerSecretAndHost(req *http.Request) (HMACSecret string, Host string, err error) {
    cfg := Config{}
    err = cfg.Get(req)
    return cfg.TrackerHMACSecret, cfg.TrackerHost, err
}
