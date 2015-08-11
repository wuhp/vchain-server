#!/bin/bash

# new_request.sh <uuid> <parent_uuid> <service> <category>

REMOTE_HOST=localhost:8010

uuid=$1
parent_uuid=$2
service=$3
category=$4

curl -X POST http://${REMOTE_HOST}/api/v1/data -H "Content-Type: application/json" -d '[
  {
    "event": "request.begin",
    "payload": {
      "uuid": "'${uuid}'",
      "parent_uuid": "'${parent_uuid}'",
      "service": "'${service}'",
      "category": "'${category}'",
      "sync_option": "sync",
      "begin_ts": '`date +%s%N`',
      "begin_metadata": {"k1":"v1", "k2":"v2"}
    }
  }
]'
