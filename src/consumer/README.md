# consumer

## Build & Install
    ./deploy.sh

## Start
    vim /vchain/server/config/consumer.json
    /vchain/server/bin/consumer -c /vchain/server/config/consumer.json

## API
### Ping
    curl -si http://localhost:8102/api/v1/ping

### Post data

    # S1-C1-001 - S2-C1-001 - S3-C1-001
    #           - S2-C2-001 - S3-C1-002

    # S1-C2-001 - S3-C2-001

    curl -si -X POST http://localhost:8102/api/v1/1/requests -d '[
        {
            "uuid": "S1-C1-001",
            "parent_uuid": "",
            "service": "S1",
            "category": "C1",
            "sync": true,
            "begin_ts": 1443965681000000000,
            "end_ts": 1443965685000000000
        },
        {
            "uuid": "S1-C2-001",
            "parent_uuid": "",
            "service": "S1",
            "category": "C2",
            "sync": true,
            "begin_ts": 1443965682000000000,
            "end_ts": 1443965685000000000
        }
    ]'

    curl -si -X POST http://localhost:8102/api/v1/1/requests -d '[
        {
            "uuid": "S2-C1-001",
            "parent_uuid": "S1-C1-001",
            "service": "S2",
            "category": "C1",
            "sync": true,
            "begin_ts": 1443965682000000000,
            "end_ts": 1443965683000000000
        },
        {
            "uuid": "S2-C2-001",
            "parent_uuid": "S1-C1-001",
            "service": "S2",
            "category": "C2",
            "sync": true,
            "begin_ts": 1443965683000000000,
            "end_ts": 1443965684000000000
        }
    ]'

    curl -si -X POST http://localhost:8102/api/v1/1/requests -d '[
        {
            "uuid": "S3-C1-001",
            "parent_uuid": "S2-C1-001",
            "service": "S3",
            "category": "C1",
            "sync": true,
            "begin_ts": 1443965682200000000,
            "end_ts": 1443965682700000000
        },
        {
            "uuid": "S3-C1-002",
            "parent_uuid": "S2-C2-001",
            "service": "S3",
            "category": "C1",
            "sync": true,
            "begin_ts": 1443965683300000000,
            "end_ts": 1443965683700000000
        },
        {
            "uuid": "S3-C2-001",
            "parent_uuid": "S1-C2-001",
            "service": "S3",
            "category": "C2",
            "sync": true,
            "begin_ts": 1443965683000000000,
            "end_ts": 1443965684000000000
        }
    ]'

    # S1-C1-001 - S2-C1-001 - S3-C1-001
    #           - S2-C2-001 - S3-C1-002

    # S1-C2-001 - S3-C2-001
    curl -si -X POST http://localhost:8102/api/v1/1/request-logs -d '[
        {
            "uuid": "S1-C1-001",
            "timestamp": 1443965681200000000,
            "log": "This is request S1-C1-001, now before S2-C1-001"
        },
        {
            "uuid": "S1-C1-001",
            "timestamp": 1443965684200000000,
            "log": "This is request S1-C1-001, now after S2-C2-001"
        },
        {
            "uuid": "S1-C2-001",
            "timestamp": 1443965682500000000,
            "log": "This is request S1-C2-001, now before S3-C2-001"
        },
        {
            "uuid": "S1-C2-001",
            "timestamp": 1443965684500000000,
            "log": "This is request S1-C2-001, now after S3-C2-001"
        }
    ]'

    curl -si -X POST http://localhost:8102/api/v1/1/request-logs -d '[
        {
            "uuid": "S2-C1-001",
            "timestamp": 1443965682100000000,
            "log": "This is request S2-C1-001, now before S3-C1-001"
        },
        {
            "uuid": "S2-C1-001",
            "timestamp": 1443965682800000000,
            "log": "This is request S2-C1-001, now after S3-C1-001"
        },
        {
            "uuid": "S2-C2-001",
            "timestamp": 1443965683200000000,
            "log": "This is request S2-C2-001, now before S3-C1-002"
        },
        {
            "uuid": "S2-C2-001",
            "timestamp": 1443965683800000000,
            "log": "This is request S2-C2-001, now after S3-C1-002"
        }
    ]'

    curl -si -X POST http://localhost:8102/api/v1/1/request-logs -d '[
        {
            "uuid": "S3-C1-001",
            "timestamp": 1443965682500000000,
            "log": "This is request S3-C1-001, in processing"
        },
        {
            "uuid": "S3-C1-002",
            "timestamp": 1443965683500000000,
            "log": "This is request S3-C1-002, in processing"
        },
        {
            "uuid": "S3-C2-001",
            "timestamp": 1443965683200000000,
            "log": "This is request S2-C2-001, in processing"
        }
    ]'

### Get data
    curl -si http://localhost:8102/api/v1/1/services
    curl -si http://localhost:8102/api/v1/1/service-chain
    curl -si http://localhost:8102/api/v1/1/request-types

    curl -si http://localhost:8102/api/v1/1/requests
    curl -si http://localhost:8102/api/v1/1/requests/S1-C1-001
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001
    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001

    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/invoke-chain
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001/invoke-chain

    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/group
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001/group

    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/root-request
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001/root-request

    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/children
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001/children

    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/parent
    curl -si http://localhost:8102/api/v1/1/requests/S1-C2-001/parent

    curl -si http://localhost:8102/api/v1/1/invoke-chains
    curl -si http://localhost:8102/api/v1/1/invoke-chains/S1/C1
    curl -si http://localhost:8102/api/v1/1/invoke-chains/S1/C1/{id}
    curl -si http://localhost:8102/api/v1/1/invoke-chains/S1/C1/{id}/root-requests

    curl -si http://localhost:8102/api/v1/1/request-logs
    curl -si http://localhost:8102/api/v1/1/requests/S3-C1-001/logs
    # curl -si http://localhost:8102/api/v1/1/request-logs/S1
    # curl -si http://localhost:8102/api/v1/1/request-logs/S1/C1
