package conf

import (
    "fmt"
    "os"
    "strings"
    "time"
)

const (
    ENV_PREFIX = "HASHTOCK_"
)

const (
    keyAPP_ADDRESS                = "APP_ADDRESS"
    keySERVE_ADDR                 = "SERVE_ADDR"
    keyDB                         = "DB"
    keyDB_NAME                    = "DB_NAME"
    keySESSION_KEY                = "SESSION_KEY"
    keySESSION_SECRET             = "SESSION_SECRET"
    keyTRACKER_URL                = "TRACKER_URL"
    keyTRACKER_SECRET             = "TRACKER_SECRET"
    keyGOOGLE_OAUTH_CLIENT_ID     = "GOOGLE_OAUTH_CLIENT_ID"
    keyGOOGLE_OAUTH_CLIENT_SECRET = "GOOGLE_OAUTH_CLIENT_SECRET"
    keyJOB_BANK_ORDERS            = "JOB_BANK_ORDERS"
    keyJOB_TAG_VALUES             = "JOB_TAG_VALUES"
    keyADMIN                      = "ADMIN"
)

var cfgHelp = map[string]string{
    keyAPP_ADDRESS:                "External app URL",
    keySERVE_ADDR:                 "Internal app host:port for binding",
    keyDB:                         "Location of MongoDB: mongodb://user:password@host:port/",
    keyDB_NAME:                    "Name of MongoDB to use",
    keySESSION_KEY:                "Session key name",
    keySESSION_SECRET:             "Session secret",
    keyTRACKER_URL:                "Tag tracker URL",
    keyTRACKER_SECRET:             "Shared secret with tracker",
    keyGOOGLE_OAUTH_CLIENT_ID:     "Google OAuth2 client id",
    keyGOOGLE_OAUTH_CLIENT_SECRET: "Google OAuth2 client secret",
    keyJOB_BANK_ORDERS:            "Time interval for running bank jobs",
    keyJOB_TAG_VALUES:             "Time interval for pulling tag values from tracker",
}

var defaultDurations = map[string]time.Duration{
    keyJOB_BANK_ORDERS: 1 * time.Minute,
    keyJOB_TAG_VALUES:  1 * time.Minute,
}

func makeKey(key string) string {
    return ENV_PREFIX + key
}

func getEnvOrDefaultDuration(key string) time.Duration {
    defaultValue, ok := defaultDurations[key]
    if !ok {
        panic(fmt.Sprintf("Default value for %v not available", key))
    }

    key = makeKey(key)
    durationStr := os.Getenv(key)
    if durationStr == "" {
        return defaultValue
    }

    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        fmt.Printf("Could not get time duration using environment key: %v. Error: %v\n", key, err)
        os.Exit(1)
    }
    return duration
}

func mustHaveValue(key string) string {
    key = makeKey(key)
    value := os.Getenv(key)
    if value == "" {
        fmt.Printf("Could not get value using environment key: %v\n\n", key)

        PrintConfHelp()

        os.Exit(1)
    }
    return value
}

func loadConfig() {
    if cfg == nil {
        cfg = new(Config)
    }

    cfg.General.AppAddress = mustHaveValue(keyAPP_ADDRESS)
    cfg.General.ServeAddr = mustHaveValue(keySERVE_ADDR)
    cfg.General.DB = mustHaveValue(keyDB)
    cfg.General.DBName = mustHaveValue(keyDB_NAME)
    cfg.General.SessionKey = mustHaveValue(keySESSION_KEY)
    cfg.General.SessionSecret = mustHaveValue(keySESSION_SECRET)

    cfg.Tracker.Url = mustHaveValue(keyTRACKER_URL)
    cfg.Tracker.HMACSecret = mustHaveValue(keyTRACKER_SECRET)

    cfg.GoogleOAuth.ClientID = mustHaveValue(keyGOOGLE_OAUTH_CLIENT_ID)
    cfg.GoogleOAuth.ClientSecret = mustHaveValue(keyGOOGLE_OAUTH_CLIENT_SECRET)

    cfg.Jobs.BankOrders = getEnvOrDefaultDuration(keyJOB_BANK_ORDERS)
    cfg.Jobs.TagValues = getEnvOrDefaultDuration(keyJOB_TAG_VALUES)

    adminKey := makeKey(keyADMIN)
    for _, e := range os.Environ() {
        keyValue := strings.Split(e, "=")
        if strings.HasPrefix(keyValue[0], adminKey) {
            cfg.General.Admin = append(cfg.General.Admin, keyValue[1])
        }
    }
}

func PrintConfHelp() {
    keysOrder := []string{
        keyAPP_ADDRESS,
        keySERVE_ADDR,
        keyDB,
        keyDB_NAME,
        keySESSION_KEY,
        keySESSION_SECRET,
        keyTRACKER_URL,
        keyTRACKER_SECRET,
        keyGOOGLE_OAUTH_CLIENT_ID,
        keyGOOGLE_OAUTH_CLIENT_SECRET,
        keyJOB_BANK_ORDERS,
        keyJOB_TAG_VALUES,
    }

    fmt.Println("Environmental variables used in configuration")
    for _, key := range keysOrder {
        help := cfgHelp[key]

        envKey := makeKey(key)
        fmt.Println(envKey)
        value := os.Getenv(envKey)
        if value == "" {
            value = "not set"
            if duration, ok := defaultDurations[key]; ok {
                value = fmt.Sprintf("%v (default)", duration)
            }
        }
        fmt.Println("\tValue:", value)
        fmt.Println("\tHelp:", help)
        fmt.Println("")
    }
}
