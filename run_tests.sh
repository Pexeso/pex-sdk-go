#!/bin/bash

export AE_SERVICE_ADDRESS=localhost:6060

curl -k -o $PWD/ca.crt "https://$AE_SERVICE_ADDRESS/roots"
export AE_SERVICE_ROOTS_FILE=$PWD/ca.crt

GRPC_VERBOSITY=debug go test
