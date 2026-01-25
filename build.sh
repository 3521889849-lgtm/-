#!/usr/bin/env bash
set -euo pipefail

RUN_NAME="piaowu"
TARGET_PKG="./cmd/dev"
GO_LDFLAGS="-checklinkname=0"

mkdir -p output/bin
cp script/* output/
chmod +x output/bootstrap.sh

if [ "${IS_SYSTEM_TEST_ENV:-0}" != "1" ]; then
    go build -ldflags="${GO_LDFLAGS}" -o output/bin/${RUN_NAME} ${TARGET_PKG}
else
    go test -c -covermode=set -ldflags="${GO_LDFLAGS}" -o output/bin/${RUN_NAME} -coverpkg=./... ${TARGET_PKG}
fi
