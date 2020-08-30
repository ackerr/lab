package utils

import "os"

// GetEnv : get env, if not exist use fallback value
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
