#!/bin/bash

export AE_SERVICE_ADDRESS=localhost:6060
#export AE_SERVICE_ROOTS_FILE=/home/stepan/Pex/ae-sdk/data/crt
export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH=/home/stepan/Pex/ae-sdk/data/crt

GRPC_VERBOSITY=debug go test
