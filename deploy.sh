#!/bin/bash

echo -e "\n+++++ Starting deployment +++++\n"

rm -rf ./bin

echo "+++++ build go packages +++++"

cd source/tweet
go get
go test ./...
env GOOS=linux GOARCH=amd64 go build -o ../../bin/tweet

echo "+++++ tweet module +++++"
cd ../deploy
go run .

echo -e "\n+++++ Deployment done +++++\n"