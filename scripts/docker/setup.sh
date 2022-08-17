#!/bin/bash

set -euo pipefail

DIR=$(dirname "$0")

CONTAINER_NAME=api

echo " > Update packages"
docker compose run --rm $CONTAINER_NAME go mod vendor
