#!/bin/bash

set -euo pipefail

protoc -I ./ ./pkg/semrel/*.proto ./pkg/provider/*.proto --go_opt=paths=source_relative --go_out=./
protoc -I ./ ./pkg/plugin/*.proto --go_opt=paths=source_relative --go_out=./ --go-grpc_opt=paths=source_relative --go-grpc_out=./
