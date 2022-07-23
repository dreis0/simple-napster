package utils

import "os"

func FromEnv(key string, fallback string) string {
	env, ok := os.LookupEnv(key)
	if !ok || env == "" {
		return fallback
	}

	return env
}
