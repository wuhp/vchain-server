# vchain-server

## components
    vchaind  - auth & project management & communicated with web browser
    receiver - communicate with vchain-client, collect end user's vchain data
    consumer - consume chain data, read by vchaind and wrote by receiver
    gateway  - providion db instance, which store end user's vchain data
    worker   - backend work deamon, parse chain data

## vchaind
    src/vchaind/README.md

## receiver
    src/receiver/README.md

## consumer
    src/consumer/README.md

## gateway
    src/gateway/README.md

## worker
    src/worker/README.md
