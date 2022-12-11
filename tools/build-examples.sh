#!/usr/bin/env bash

set -e

pushd "$(dirname "${BASH_SOURCE[0]}")/../examples"

for d in */; do
  pushd "${d}" >/dev/null
  echo "Building ${d}"
  go build
  popd >/dev/null
done