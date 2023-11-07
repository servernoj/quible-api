package main

import (
	"os"

	"github.com/joho/godotenv"
)

func InitializeEnvironment(env string) error {
	fn := ".env"
	if len(env) > 0 {
		fn += "." + env
	}
	if _, err := os.Stat(fn); err == nil {
		return godotenv.Load(fn)
	}
	return nil
}
