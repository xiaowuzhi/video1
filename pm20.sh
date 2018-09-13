#! /bin/bash

# Build web and other services
cd ~/awesomeProject/go/bin


pm2 delete api scheduler streamserver web
pm2 start api scheduler streamserver web
