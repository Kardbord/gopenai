#!/usr/bin/env bash

set -e

pushd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null

./examples-build.sh

pushd .. >/dev/null
echo "Formatting code..."
go fmt ./...
echo "Running tests..."
go test ./...
echo "Building $(basename "$(pwd)")..."
go build ./...
echo "Vetting..."
go vet ./...
echo "Done."
