package config

import (
	"fmt"
	"os"
)

func MustExistsEnv(key string) error {
	k := os.Getenv(key)
	if len(k) == 0 {
		return fmt.Errorf("env var  %s was not found", k)
	}
	return nil
}
