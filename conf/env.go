package conf

import (
	"fmt"
	"os"
	"time"
)

const (
	ENV_PREFIX = "HASHTOCK_"
)

const (
	keyAuthAddress     = "AUTH_ADDRESS"
	keySERVE_ADDR      = "SERVE_ADDR"
	keyDB              = "DB"
	keyDB_NAME         = "DB_NAME"
	keyJOB_BANK_ORDERS = "JOB_BANK_ORDERS"
	keyJOB_TAG_VALUES  = "JOB_TAG_VALUES"
	keyNATS            = "NATS"
)

var cfgHelp = map[string]string{
	keyAuthAddress:     "Host and port for the auth service",
	keySERVE_ADDR:      "Internal app host:port for binding",
	keyDB:              "Location of MongoDB: mongodb://user:password@host:port/",
	keyDB_NAME:         "Name of MongoDB to use",
	keyJOB_BANK_ORDERS: "Time interval for running bank jobs",
	keyJOB_TAG_VALUES:  "Time interval for pulling tag values from tracker",
	keyNATS:            "Location of natds server",
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

	cfg.General.AuthAddress = mustHaveValue(keyAuthAddress)
	cfg.General.ServeAddr = mustHaveValue(keySERVE_ADDR)
	cfg.General.DB = mustHaveValue(keyDB)
	cfg.General.DBName = mustHaveValue(keyDB_NAME)
	cfg.General.NATS = mustHaveValue(keyNATS)

	cfg.Jobs.BankOrders = getEnvOrDefaultDuration(keyJOB_BANK_ORDERS)
	cfg.Jobs.TagValues = getEnvOrDefaultDuration(keyJOB_TAG_VALUES)
}

func PrintConfHelp() {
	keysOrder := []string{
		keyAuthAddress,
		keySERVE_ADDR,
		keyDB,
		keyDB_NAME,
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
