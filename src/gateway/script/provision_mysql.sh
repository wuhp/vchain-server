#!/bin/bash

set -ex

export DATABASE_HOST=$1
export DATABASE_PORT=$2
export DATABASE_USER=$3
export DATABASE_PASSWD=$4
export DATABASE_DBNAME=$5

export MIGRATION_SCRIPT_DIR="/vchain/server/datasource/schema_migration"

mysql -h ${DATABASE_HOST} -P ${DATABASE_PORT} -u${DATABASE_USER} -p${DATABASE_PASSWD} \
      -e "CREATE DATABASE \`${DATABASE_DBNAME}\` CHARACTER SET utf8 COLLATE utf8_general_ci;"

/vchain/server/bin/migrate latest
