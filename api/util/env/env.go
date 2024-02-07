package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var envVars = make(map[string]string)

func init() {
	parseEnvVars()
}

func parseEnvVars() {
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			envVars[e[:i]] = e[i+1:]
		}
	}

	env := GetEnv("ENV")
	if env == "" {
		env = "local"
		os.Setenv("ENV", "local")
		envVars["ENV"] = "local"
	}

	// Load from .env files
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

// GetEnv returns the environment variable for the given key (variable name)
func GetEnv(key string) string {
	if val, ok := envVars[key]; ok {
		return val
	}
	return ""
}
