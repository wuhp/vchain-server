# gateway

## Build & Install
    ./deploy.sh

## SQL Migration
    vim /vchain/migration/gateway.def 
    . /vchain/migration/gateway.source
    /vchain/migration/db_create
    /vchain/migration/migrate latest

## Start
    vim /vchain/server/config/gateway.json
    /vchain/server/bin/gateway -c /vchain/server/config/gateway.json

## API
### Ping
    curl -si http://localhost:8101/api/v1/ping

### Mysql server
    curl -si http://localhost:8101/api/v1/servers
    curl -si -X POST http://localhost:8101/api/v1/servers -d '{
        "host": "192.168.38.132",
        "port": 13306,
        "admin_user": "root",
        "admin_password": "password"
    }'
    curl -si http://localhost:8101/api/v1/servers/1
    curl -si -X PUT http://localhost:8101/api/v1/servers/1 -d '{
        "admin_user": "root",
        "admin_password": "root"
    }'
    curl -si -X DELETE http://localhost:8101/api/v1/servers/1

### Mysql instance
    curl -si http://localhost:8101/api/v1/mysql
    curl -si -X POST http://localhost:8101/api/v1/provision -d '{"project_id":1}'
    curl -si http://localhost:8101/api/v1/mysql/1
    curl -si -X POST http://localhost:8101/api/v1/unprovision -d '{"project_id":1}'

