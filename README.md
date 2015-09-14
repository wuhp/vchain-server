# vchain-server

## SQL migration
    mysql -h localhost -P 3306 -uroot -proot vchain < migration/1/forward.sql

## Build
    export GOPATH=$PWD
    go get account
    go install account

## Start
    ./bin/account -c config/account.json
