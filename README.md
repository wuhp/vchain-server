# vchain-server

## DB
    mysql -h 192.168.38.132 -P 3306 -uroot -proot -e "CREATE DATABASE \`vchaind\` CHARACTER SET utf8 COLLATE utf8_general_ci;"

## SQL migration
    mysql -h 192.168.38.132 -P 3306 -uroot -proot vchaind < src/vchaind/migration/1/forward.sql

## Build
    export GOPATH=$PWD
    go get vchaind
    go install vchaind

## Start
    ./bin/vchaind -c config/vchaind.json
