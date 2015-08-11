#!/bin/bash

# S1-C1-001 - S2-C1-001 - S3-C1-001
#           - S2-C2-001 - S3-C1-002

# S1-C2-001 - S3-C2-001

`dirname $0`/tools/new_request.sh "S1-C1-001" "" "S1" "C1"
`dirname $0`/tools/new_request.sh "S2-C1-001" "S1-C1-001" "S2" "C1"
`dirname $0`/tools/new_request.sh "S3-C1-001" "S2-C1-001" "S3" "C1"
`dirname $0`/tools/finish_request.sh "S3-C1-001"
`dirname $0`/tools/finish_request.sh "S2-C1-001"
`dirname $0`/tools/new_request.sh "S2-C2-001" "S1-C1-001" "S2" "C2"
`dirname $0`/tools/new_request.sh "S3-C1-002" "S2-C2-001" "S3" "C1"
`dirname $0`/tools/finish_request.sh "S3-C1-002"
`dirname $0`/tools/finish_request.sh "S2-C2-001"
`dirname $0`/tools/finish_request.sh "S1-C1-001"

`dirname $0`/tools/new_request.sh "S1-C2-001" "" "S1" "C2"
`dirname $0`/tools/new_request.sh "S3-C2-001" "S1-C2-001" "S3" "C2"
`dirname $0`/tools/finish_request.sh "S3-C2-001"
`dirname $0`/tools/finish_request.sh "S1-C2-001"
