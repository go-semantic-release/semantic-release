#!/usr/bin/env bash

set -euo pipefail

for f in $(find ./pkg -name "*.proto") ; do
  echo "generating $f"
  protoc -I ./ $f --go_opt=paths=source_relative --go_out=./ --go-grpc_opt=paths=source_relative --go-grpc_out=./
done
