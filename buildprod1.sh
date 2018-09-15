#! /bin/bash

# Build web and other services
cd ~/awesomeProject/go/src/video1/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ~/awesomeProject/go/src/video1/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ~/awesomeProject/go/src/video1/streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd ~/awesomeProject/go/src/video1/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web
