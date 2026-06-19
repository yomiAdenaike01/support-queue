package images

import (
	"os"
)

func GetEnvVar(varName string, defeaultVal string) string {
	envVar := os.Getenv(varName)
	if envVar == "" {
		return defeaultVal
	}
	return envVar
}
