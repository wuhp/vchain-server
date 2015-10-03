#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))/../..

# Build
export GOPATH=${ROOT_DIR}
go get worker
go install worker

# Install
mkdir -p /vchain/server/bin
install ${ROOT_DIR}/bin/worker /vchain/server/bin/worker
