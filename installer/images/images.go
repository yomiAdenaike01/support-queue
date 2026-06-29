package images

import (
	"os"
)

func getEnvVar(varName string, defeaultVal string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defeaultVal
	}
	return envVar
}

func GetruntimeVersion() string {
	return getEnvVar("VERSION", "0.0.1-dev")
}
