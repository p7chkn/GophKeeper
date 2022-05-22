#!/bin/sh

clear

TRG_PKG='main'

AppVersion=N/A

BUILD_TIME=$(date +"%Y-%m-%dT%H:%M:%S%z")

GV=$(git tag || echo 'N/A')
# shellcheck disable=SC2039
if [[ $GV =~ [^[:space:]]+ ]];
then
    AppVersion=${BASH_REMATCH[0]}
fi

AppVersion=$(git tag --sort=-version:refname | head -n 1)

FLAG="-X $TRG_PKG.BuildTime=$BUILD_TIME"
FLAG="$FLAG -X $TRG_PKG.AppVersion=$AppVersion"

GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe -ldflags "$FLAG" main.go

GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin -ldflags "$FLAG" main.go

GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux -ldflags "$FLAG" main.go

