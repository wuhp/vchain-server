#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))/../..

# Build
export GOPATH=${ROOT_DIR}
go get consumer
go install consumer

# Install
mkdir -p /vchain/server/bin
install ${ROOT_DIR}/bin/consumer /vchain/server/bin/consumer
