# receiver

## Build & Install
    ./deploy.sh

## Start
    vim /vchain/server/config/receiver.json
    /vchain/server/bin/receiver -c /vchain/server/config/receiver.json

## API
### Ping
    curl -si http://localhost:8100/api/v1/ping

### Post data
    curl -si http://localhost:8100/api/v1/data -X POST -H "Authorization:xyz" -d '{
      "requests":[
        {"uuid":"uuid003","parent_uuid":"uuid002","service":"S3","category":"C1","sync":true,"begin_ts":1444065727501828157,"end_ts":1444065727501842036},
        {"uuid":"uuid002","parent_uuid":"uuid001","service":"S2","category":"C1","sync":true,"begin_ts":1444065727501818230,"end_ts":1444065727501876188}
      ],
      "request-logs":[
        {"uuid":"uuid001","timestamp":1444065727501668442,"log":"r1 uuid001, before calling r2"},
        {"uuid":"uuid002","timestamp":1444065727501818528,"log":"r2 uuid002, before calling r3"},
        {"uuid":"uuid003","timestamp":1444065727501828438,"log":"r3 uuid003"},
        {"uuid":"uuid002","timestamp":1444065727501870527,"log":"r2 uuid002, after calling r3"}
      ]
    }'
