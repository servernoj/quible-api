#!/usr/bin/env bash

ROOT=$(dirname ${0})
CMD=$1

if [ ! -d $ROOT/$CMD ]; then 
  >&2 echo "Unable to locate command '$CMD'"
  exit 1
fi

cd $ROOT
go mod download

go run ./$CMD "${@:2}"





