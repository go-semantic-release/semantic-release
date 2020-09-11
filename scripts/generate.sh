#!/usr/bin/env bash

set -euo pipefail

if [[ "$OSTYPE" != "darwin"* ]]; then
  shopt -s globstar
fi

for f in $(ls ./pkg/**/*.proto) ; do
  echo "generating $f"
  protoc -I ./ $f --go_opt=paths=source_relative --go_out=./ --go-grpc_opt=paths=source_relative --go-grpc_out=./
done
