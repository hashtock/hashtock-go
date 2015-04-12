package conf

import (
    "fmt"
    "net/http"
    "time"

    "golang.org/x/oauth2"
)

type Config struct {
    General GeneralConf

    Tracker struct {
        HMACSecret string
        Url        string
    }

    GoogleOAuth struct {
        ClientID     string
        ClientSecret string
    }

    Jobs struct {
        BankOrders time.Duration
        TagValues  time.Duration
    }
}

type GeneralConf struct {
    AppAddress    string
    ServeAddr     string
    DB            string
    DBName        string
    SessionKey    string
    SessionSecret string
    Admin         []string

    admins map[string]bool
}

func (c *Config) isAdmin(email string) bool {
    if c.General.admins == nil {
        c.General.admins = make(map[string]bool, len(c.General.Admin))
        for _, admin := range c.General.Admin {
            c.General.admins[admin] = true
        }
    }

    return c.General.admins[email]
}

var cfg *Config = nil

const exampleConfig = `[general]
AppAddress = "www.name.com"
ServeAddr = ":8080"
DB = "localhost"
DBName = "Hashtock"
SessionKey = "session-key-name"
SessionSecret = "this-is-a-secter"
Admin = "me@here.com"
Admin = "other@here.com"

[tracker]
URL = "www.tracker.com:80"
HMACSecret = "shared secret with tracker"

[GoogleOAuth]
ClientID = "ID of Google auth key"
ClientSecret = "Shared secret"

[Jobs]
BankOrders = 1m
TagValues = 1m
`

func GetConfig() *Config {
    if cfg == nil {
        loadConfig()
    }

    return cfg
}

func TrackerSecretAndHost(req *http.Request) (HMACSecret string, Url string, err error) {
    cfg := GetConfig()
    return cfg.Tracker.HMACSecret, cfg.Tracker.Url, nil
}

func (c *Config) GAuthConfig() oauth2.Config {
    oauth2Redirect := fmt.Sprintf("%v/auth/oauth2callback", c.General.AppAddress)
    return oauth2.Config{
        ClientID:     c.GoogleOAuth.ClientID,
        ClientSecret: c.GoogleOAuth.ClientSecret,
        Scopes: []string{
            "profile",
            "https://www.googleapis.com/auth/userinfo.profile",
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/plus.login",
            "https://www.googleapis.com/auth/plus.me",
        },
        RedirectURL: oauth2Redirect,
    }
}

func IsAdmin(email string) bool {
    return GetConfig().isAdmin(email)
}
