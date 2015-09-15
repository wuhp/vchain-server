# vchain-server

## DB
    mysql -h 192.168.38.132 -P 3306 -uroot -proot -e "CREATE DATABASE \`account\` CHARACTER SET utf8 COLLATE utf8_general_ci;"

## SQL migration
    mysql -h 192.168.38.132 -P 3306 -uroot -proot account < migration/1/forward.sql

## Build
    export GOPATH=$PWD
    go get account
    go install account

## Start
    ./bin/account -c config/account.json
