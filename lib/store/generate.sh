#!/usr/bin/env bash


SQL_BOILER_BIN="sqlboiler"
SQL_BOILER_DRIVER="psql"
SQL_BOILER_DRIVER_UPPERCASE=$(tr '[:lower:]' '[:upper:]' <<< "$SQL_BOILER_DRIVER")
SQL_BOILER_DBNAME="${SQL_BOILER_DRIVER_UPPERCASE}_DBNAME"
SQL_BOILER_USER="${SQL_BOILER_DRIVER_UPPERCASE}_USER"
SQL_BOILER_PASS="${SQL_BOILER_DRIVER_UPPERCASE}_PASS"

ENV_FILE_NAME="../../.env"
ENV_FILE_BASE=$(basename "${ENV_FILE_NAME}")
ENV_FILE_DIR=$(dirname "${ENV_FILE_NAME}")
ENV_FILE_ABS=$(cd "${ENV_FILE_DIR}";pwd)/${ENV_FILE_BASE}

if ! [ -f ${ENV_FILE} ]; then
  >&2 echo "Unable to locate .env file at `${ENV_FILE_ABS}`"
  exit 1
fi

if ! which ${SQL_BOILER_BIN} > /dev/null 2> /dev/null; then
  >&2 echo "Unable to locate SQLBoiler binary"
fi

source ${ENV_FILE_NAME}

eval export ${SQL_BOILER_DBNAME}="${POSTGRES_DB}"
eval export ${SQL_BOILER_USER}="${POSTGRES_USER}"
eval export ${SQL_BOILER_PASS}="${POSTGRES_PASSWORD}"

${SQL_BOILER_BIN} --wipe --add-global-variants ${SQL_BOILER_DRIVER}





