package env

import (
	"fmt"
	"os"
	"strconv"
	"log"
)

var (
	Port = getEnvInt("SERVER_HTTP_PORT", 9194)
)

func getEnvInt(key string, fallback int) int {
	k, ok := os.LookupEnv(key)

	if ok && k != "" {
		val, err := strconv.Atoi(k)
		if err != nil {
			errMsg := fmt.Errorf("env var key '%s': '%s' defaulting to '%d' error: %s", key, k, fallback, err.Error())
			log.Error("failed to parse env var using default", slog.Any("error", errMsg))
			return fallback
		}
		return val
	}

	return fallback
}

func getEnvStr(key string, fallback string) string {
	k, ok := os.LookupEnv(key)
	if ok && k != "" {
		return k
	}

	return fallback
}

func getEnvBoolean(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)

	if !ok || value == "" {
		return fallback
	}

	booleanValue, error := strconv.ParseBool(value)
	if error != nil {
		return fallback
	}

	return booleanValue
}
