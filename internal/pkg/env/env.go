package env

import (
	"fmt"
	"os"
)

func CheckEnvVars(vars []string) error {
	missingVars := []string{}
	for _, v := range vars {
		if value, exists := os.LookupEnv(v); !exists || value == "" {
			missingVars = append(missingVars, v)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("you should provide env vars: %v", missingVars)
	}
	return nil
}
