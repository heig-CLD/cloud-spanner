#!/bin/bash

# https://cloud.google.com/spanner/docs/emulator#docker
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set project noted-episode-316407
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

# Creating an instance
gcloud spanner instances create test-instance --config=emulator --description="Test instance" --nodes=1
gcloud spanner databases create test-database --instance test-instance
