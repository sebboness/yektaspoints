#!/usr/bin/env bash
set -e

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

CUR_DIR=$(pwd)
VERSION=$(cat ./VERSION)
TZ=UTC
BUILT_AT=$(date +%FT%$TZ)

export $(cat $CUR_DIR/.env.local | xargs) && env

echo "VERSION=" $VERSION
echo "TZ=" $TZ
echo "BUILT_AT=" $BUILT_AT
echo "APPNAME=" $APPNAME
echo "RUN_AS_WEB_API=" $RUN_AS_WEB_API

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

cd $BASE_DIR/cmd/lambda
go build -tags lambda.norpc -o ./webapi -ldflags "-s -w -X \"main.Version=$VERSION\" -X \"main.BuiltAt=$BUILT_AT\""
go run .