# consumer

## Build & Install
    ./deploy.sh

## Start
    vim /vchain/server/config/consumer.json
    /vchain/server/bin/consumer -c /vchain/server/config/consumer.json

## API
### Ping
    curl -si http://localhost:8101/api/v1/ping
