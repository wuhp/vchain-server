# vchain-server

## SQL migration
mysql -h 192.168.38.128 -P 4306 -uroot -proot vchaind < migration/1/forward.sql

## Build
export GOPATH=$PWD
go install vchaind

## Start
./bin/vchaind -c example/vchaind.json

## Examples
### Import data
./example/example001.sh

### Fetch data
curl http://localhost:8010/api/v1/services

curl http://localhost:8010/api/v1/services-chain

curl http://localhost:8010/api/v1/services/S1/request-categories

curl http://localhost:8010/api/v1/invoke-chains

curl http://localhost:8010/api/v1/invoke-chains/S1/C2

curl http://localhost:8010/api/v1/invoke-chains/S1/C2/2

curl http://localhost:8010/api/v1/invoke-chains/S1/C2/2/root-requests

curl http://localhost:8010/api/v1/request-overview

curl http://localhost:8010/api/v1/request-types

curl http://localhost:8010/api/v1/requests

curl http://localhost:8010/api/v1/requests/S1-C1-001

curl http://localhost:8010/api/v1/requests/S1-C1-001/invoke-chain

curl http://localhost:8010/api/v1/requests/S1-C1-001/request-group

curl http://localhost:8010/api/v1/requests/S1-C1-001/children

curl http://localhost:8010/api/v1/requests/S2-C1-001/parent

curl http://localhost:8010/api/v1/requests/S2-C1-001/root-request
