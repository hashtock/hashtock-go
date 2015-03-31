package conf

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "code.google.com/p/gcfg"
    "golang.org/x/oauth2"
)

type confDuration string

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
        BankOrders confDuration
        TagValues  confDuration
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

func (c confDuration) Duration() time.Duration {
    duration, err := time.ParseDuration(string(c))
    if err != nil {
        log.Fatalln("Could not parse duraiton:", err)
    }
    return duration
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

func loadConfig() {
    if cfg == nil {
        cfg = new(Config)
    }
    err := gcfg.ReadFileInto(cfg, "config.ini")
    if err != nil {
        if os.IsNotExist(err) {
            log.Fatalf("Could not find remote tracker configuration. Expected config.ini with content:\n%v\n", exampleConfig)
        } else {
            log.Fatalln("Config error:", err.Error())
        }
    }
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
