#!/bin/sh

# Use port 9010 for gRPC requests.
# Use port 9020 for REST requests.
docker run -p 9010:9010 -p 9020:9020 gcr.io/cloud-spanner-emulator/emulator
