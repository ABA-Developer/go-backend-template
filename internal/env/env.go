package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}

	return valInt
}
func GetBool(key string, fallback bool) bool {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valBool, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}

	return valBool
}
func GetFloat(key string, fallback float64) float64 {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valFloat, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fallback
	}

	return valFloat
}
