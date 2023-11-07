package main

import (
	"github.com/joho/godotenv"
)

func InitializeEnvironment(env string) error {
	fn := ".env"
	if len(env) > 0 {
		fn += "." + env
	}
	return godotenv.Load(fn)
}
