#!/bin/bash

set -uo pipefail

DIR=$(dirname "$0")

set -e

echo " > Setup app"
"${DIR}"/setup.sh
