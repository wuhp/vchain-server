#!/bin/bash

set -ex

export DATABASE_HOST=$1
export DATABASE_PORT=$2
export DATABASE_USER=$3
export DATABASE_PASSWD=$4
export DATABASE_DBNAME=$5

mysql -h ${DATABASE_HOST} -P ${DATABASE_PORT} -u${DATABASE_USER} -p${DATABASE_PASSWD} \
      -e "DROP DATABASE \`${DATABASE_DBNAME}\`;"
