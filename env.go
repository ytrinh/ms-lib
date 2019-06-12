package ms

import (
	"os"
)

// GetEnv return an env value or defaultValue if env value does not exist
func GetEnv(envName string, defaultValue string) string {
    value := os.Getenv(envName)
    if len(value) > 0 {
        return value
    }

    return defaultValue
}
