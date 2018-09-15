#! /bin/bash

# Build web and other services
cd ./bin
pm2 restart api scheduler streamserver web
