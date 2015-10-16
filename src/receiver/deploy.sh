#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))/../..

# Build
export GOPATH=${ROOT_DIR}
go get receiver
go install receiver

# Install
mkdir -p /vchain/server/bin
install ${ROOT_DIR}/bin/receiver /vchain/server/bin/receiver
