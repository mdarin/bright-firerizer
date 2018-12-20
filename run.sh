#!/bin/bash

echo "prepareing evn..."
export GOPATH=$(pwd)
echo "starting the server..."
go run ./src/main.go 
exit 0
