#! /bin/bash

# Build web and other services
cd ~/work/src/video1/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ~/work/src/video1/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ~/work/src/video1/streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd ~/work/src/video1/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web