# gateway

## SQL migration
    mysql -h 192.168.38.132 -P 3306 -uroot -proot -e "CREATE DATABASE \`gateway\` CHARACTER SET utf8 COLLATE utf8_general_ci;"
    mysql -h 192.168.38.132 -P 3306 -uroot -proot gateway < migration/1/forward.sql

## Build
    export GOPATH=$PWD
    go get gateway
    go install gateway

## Start
    ./bin/gateway -c config/gateway.json
