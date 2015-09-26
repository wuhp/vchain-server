#!/bin/bash

# Only support mysql

usage() {
  echo "Usage: "
  echo "  $0 ping"
  echo "  $0 current"
  echo "  $0 latest"
  echo "  $0 <migration num>"
}

ping_mysql() {
  echo quit | ${MYSQL_CLIENT} > /dev/null 2>&1
  return $?
}

show_cur_version() {
  echo "show tables;" | ${MYSQL_CLIENT} | grep "schema_version" > /dev/null 2>&1
  if [ $? -eq 0 ]; then
    echo "select version from schema_version;" | ${MYSQL_CLIENT} | tail -n 1
    return
  fi

  echo "create table schema_version(version int);" | ${MYSQL_CLIENT} > /dev/null 2>&1
  echo "insert into schema_version(version) values(0);" | ${MYSQL_CLIENT} > /dev/null 2>&1
  echo 0
}

forward() {
  echo "Start schema migration, from $1 forward to $2 ..."
  for n in `seq $(($1 + 1)) $2`
  do
    ${MYSQL_CLIENT} < ${MIGRATION_SCRIPT_DIR}/${n}/forward.sql
    if [ $? -eq 0 ]; then
      echo "Forward migration on ${n}   [   OK   ]"
      echo "update schema_version set version=${n};" | ${MYSQL_CLIENT}
      continue
    fi
    echo "Forward migration on ${n}   [ FAILED ]"
    return -1
  done
}

backward() {
  echo "Start schema migration, from $1 backward to $2 ..."
  for n in `seq $1 -1 $(($2 + 1))`
  do
    ${MYSQL_CLIENT} < ${MIGRATION_SCRIPT_DIR}/${n}/backward.sql
    if [ $? -eq 0 ]; then
      echo "Backward migration on ${n}   [   OK   ]"
      echo "update schema_version set version=${n}-1;" | ${MYSQL_CLIENT}
      continue
    fi
    echo "Backward migration on ${n}   [ FAILED ]"
    return -1
  done
}

[ $# -ne 1 ] && usage && exit -1
[ "$1" = "help" -o "$1" = "-h" -o "$1" = "--help" ] && usage && exit 0

[ -z ${DATABASE_HOST} ] && echo "ERROR: env DATABASE_HOST not set" && exit -1
[ -z ${DATABASE_PORT} ] && echo "ERROR: env DATABASE_PORT not set" && exit -1
[ -z ${DATABASE_DBNAME} ] && echo "ERROR: env DATABASE_DBNAME not set" && exit -1
[ -z ${DATABASE_USER} ] && echo "ERROR: env DATABASE_USER not set" && exit -1
[ -z ${DATABASE_PASSWD} ] && echo "ERROR: env DATABASE_PASSWD not set" && exit -1

[ ! -d ${MIGRATION_SCRIPT_DIR} ] && echo "ERROR: value of MIGRATION_SCRIPT_DIR not exist" && exit -1

which mysql > /dev/null 2>&1
[ $? -ne 0 ] && echo "ERROR: can not find mysql client" && exit -1

MYSQL_CLIENT="mysql -h${DATABASE_HOST} -P${DATABASE_PORT} -u${DATABASE_USER} -p${DATABASE_PASSWD} -D${DATABASE_DBNAME}"

### Ping ###

if [ "$1" = "ping" ]; then
  ping_mysql
  if [ $? -ne 0 ]; then
    echo "ERROR: can not connect to mysql instance"
    exit -1
  else
    echo "OK"
    exit 0
  fi
fi

### Show current version ###

if [ "$1" = "current" ]; then
  ping_mysql
  [ $? -ne 0 ] && echo "ERROR: can not connect to mysql instance" && exit -1
  show_cur_version
  exit 0
fi

### Migration ###

ping_mysql
[ $? -ne 0 ] && echo "ERROR: can not connect to mysql instance" && exit -1

current_version=`show_cur_version`
dest_version=$1

if [ "$1" = "latest" ]; then
  dest_version=`ls ${MIGRATION_SCRIPT_DIR} | grep "^[0-9][0-9]*" | sort -g | tail -n 1`
  [ -z ${dest_version} ] && echo "ERROR: can not find latest migration version" && exit -1
fi

[ ! -d ${MIGRATION_SCRIPT_DIR}/${dest_version} ] && echo "ERROR: can not find migration num ${dest_version}" && exit -1

if [ ${dest_version} -gt ${current_version} ]; then
  forward ${current_version} ${dest_version}
  exit $?
fi

if [ ${dest_version} -lt ${current_version} ]; then
  backward ${current_version} ${dest_version}
  exit $?
fi

exit 0
