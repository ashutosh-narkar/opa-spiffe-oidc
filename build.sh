#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

(cd src/backend-app && GOOS=linux go build -v -o $DIR/docker/backend/backend-app)
(cd src/invoice-app && GOOS=linux go build -v -o $DIR/docker/invoice/invoice-app)
