#!/usr/bin/env bash

if [ "--${ENV_DSN}--" = "----" ]; then
  >&2 echo "Environment variable `ENV_DSN` is not set"
  exit 1
fi

cd $(dirname ${0})

go mod download
go run ./goose.go "$@"





