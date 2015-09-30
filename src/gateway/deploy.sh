#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))/../..

# Build
export GOPATH=${ROOT_DIR}
go get gateway
go install gateway

# Install
mkdir -p /vchain/server/bin
mkdir -p /vchain/server/datasource

install ${ROOT_DIR}/bin/gateway /vchain/server/bin/gateway
install ${ROOT_DIR}/tools/migrate /vchain/server/bin/migrate
install ${ROOT_DIR}/src/gateway/script/migrate_mysql.sh /vchain/server/bin/migrate_mysql.sh
install ${ROOT_DIR}/src/gateway/script/provision_mysql.sh /vchain/server/bin/provision_mysql.sh
install ${ROOT_DIR}/src/gateway/script/unprovision_mysql.sh /vchain/server/bin/unprovision_mysql.sh

rm -rf /vchain/server/datasource/schema_migration
cp -r ${ROOT_DIR}/src/datasource/schema_migration /vchain/server/datasource/schema_migration

# SQL Migration
mkdir -p /vchain/migration

install ${ROOT_DIR}/tools/migrate /vchain/migration/migrate
install ${ROOT_DIR}/tools/db_create /vchain/migration/db_create
install ${ROOT_DIR}/tools/db_destroy /vchain/migration/db_destroy

rm -rf /vchain/migration/gateway
cp -r ${ROOT_DIR}/src/gateway/migration /vchain/migration/gateway
