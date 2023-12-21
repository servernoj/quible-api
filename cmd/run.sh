#!/usr/bin/env bash

# Usage: ./run.sh <command> <args>
# with <command> being a sub-directory of `cmd` with compilable `main.go`
#      <args> being command specific arguments, see below example
#
# For example, to migrate DB to the latest version run `./run.sh migrate up`

ROOT=$(dirname ${0})
CMD=$1

if [ ! -d $ROOT/$CMD ]; then 
  >&2 echo "Unable to locate command '$CMD'"
  exit 1
fi

cd $ROOT
go mod download

go run ./$CMD "${@:2}"





