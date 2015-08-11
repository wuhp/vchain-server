#!/bin/bash

# finish_request.sh <uuid>

REMOTE_HOST=localhost:8010

uuid=$1

curl -X POST http://${REMOTE_HOST}/api/v1/data -H "Content-Type: application/json" -d '[
  {
    "event": "request.end",
    "payload": {
      "uuid": "'${uuid}'",
      "end_ts": '`date +%s%N`',
      "end_metadata": {"k3":"v3", "k4":"v4"}
    }
  }
]'
